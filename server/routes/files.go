package routes

import (
	"encoding/json"
	"github.com/akrantz01/bookpi/server/responses"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Routes for file management
func Files(filesDirectory string, router *mux.Router) {
	router.PathPrefix("/files").HandlerFunc(fileRouter(filesDirectory))
}

// Handle routing based on methods for files
func fileRouter(filesDirectory string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Assemble full path
		p := path.Join(filesDirectory, r.Header.Get("X-BPI-Username"), strings.TrimPrefix(r.URL.Path, "/api/files"))

		switch r.Method {
		case http.MethodGet:
			listFiles(w, r, p)

		case http.MethodPost:
			createFile(w, r, p)

		case http.MethodPut:
			updateFile(w, r, p, filesDirectory)

		case http.MethodDelete:
			deleteFile(w, r, p)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// List all files in a directory or a file's information, or download a file
func listFiles(w http.ResponseWriter, r *http.Request, path string) {
	// Get file statistics
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		responses.Error(w, http.StatusNotFound, "specified file/directory does not exist")
		return
	} else if err != nil {
		log.Printf("ERROR: failed to stat file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	// Return file info if file
	if !info.IsDir() {
		// Download file if query param
		if r.URL.Query().Get("download") != "" {
			http.ServeFile(w, r, path)
			return
		}

		responses.SuccessWithData(w, map[string]interface{}{
			"name": info.Name(),
			"size": info.Size(),
			"last_modified": info.ModTime().Unix(),
			"directory": info.IsDir(),
		})
		return
	}

	// Get all files in directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("ERROR: failed to list files in directory: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to list files")
		return
	}

	// Format file info objects
	var children []map[string]interface{}
	for _, file := range files {
		children = append(children, map[string]interface{}{
			"name": file.Name(),
			"size": file.Size(),
			"last_modified": file.ModTime().Unix(),
			"directory": file.IsDir(),
		})
	}

	// Set to empty array if length zero
	if len(children) == 0 {
		children = []map[string]interface{}{}
	}

	rawPath := filepath.Clean(strings.TrimPrefix(r.RequestURI, "/api/files"))

	responses.SuccessWithData(w, map[string]interface{}{
		"name": info.Name(),
		"size": info.Size(),
		"last_modified": info.ModTime().Unix(),
		"directory": info.IsDir(),
		"root": rawPath == "/" || rawPath == ".",
		"children": children,
	})
}

// Upload a new file
func createFile(w http.ResponseWriter, r *http.Request, path string) {
	// Validate initial headers
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'multipart/form-data'")
		return
	}

	// Check if creating directory
	directory := r.URL.Query().Get("directory") != ""

	// Get file statistics
	info, err := os.Stat(path)
	if os.IsNotExist(err) && !directory {
		responses.Error(w, http.StatusNotFound, "specified directory does not exist")
		return
	} else if !os.IsNotExist(err) && directory {
		responses.Error(w, http.StatusConflict, "specified directory already exists")
		return
	} else if err != nil && ((os.IsExist(err) && directory) || (os.IsNotExist(err) && !directory)) {
		log.Printf("ERROR: failed to stat file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	// Ensure not uploading to file
	if !directory && !info.IsDir() {
		responses.Error(w, http.StatusBadRequest, "cannot upload to file")
		return
	}

	// Check if uploading directory
	if directory {
		if err := os.Mkdir(path, os.ModeDir|0755); err != nil {
			log.Printf("ERROR: fialed to create new directory under user: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to create directory")
			return
		}

		responses.Success(w)
		return
	}

	// Allow 32Mb internal buffer for upload
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("ERROR: failed to parse multipart form for file upload: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to parse form")
		return
	}

	// Get file from upload
	in, handler, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		responses.Error(w, http.StatusBadRequest, "field 'file' must be a file")
		return
	} else if err != nil {
		log.Printf("ERROR: failed to parse file from form: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to parse form")
		return
	}
	defer func() {
		if err := in.Close(); err != nil {
			log.Printf("ERROR: falied to close uploaded file stream: %v\n", err)
		}
	}()

	// Check file doesn't already exist
	if _, err := os.Stat(path+"/"+handler.Filename); err == nil {
		responses.Error(w, http.StatusConflict, "file already exists")
		return
	} else if !os.IsNotExist(err) {
		log.Printf("ERROR: failed to stat output file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to check output file")
		return
	}

	// Open output file
	out, err := os.OpenFile(path+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("ERROR: failed to open output file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Printf("ERROR: failed to close output file: %v\n", err)
		}
	}()

	// Copy uploaded to output
	if _, err := io.Copy(out, in); err != nil {
		log.Printf("ERROR: failed to copy uploaded file to output file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to copy file")
		return
	}

	responses.Success(w)
}

// Change a file's name on disk
func updateFile(w http.ResponseWriter, r *http.Request, path, filesDirectory string) {
	// Don't allow changes to user root
	rawPath := filepath.Clean(strings.TrimPrefix(r.RequestURI, "/api/files"))
	if rawPath == "." || rawPath == "/" {
		responses.Error(w, http.StatusForbidden, "not allowed to move user root")
		return
	}

	// Validate initial request on headers and body existence
	if r.Header.Get("Content-Type") != "application/json" {
		responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		responses.Error(w, http.StatusBadRequest, "request body must be present")
		return
	}

	// Parse and validate body fields
	var body struct{
		Filename string `json:"filename"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responses.Error(w, http.StatusBadRequest, "invalid json format for request body")
		return
	} else if body.Path != "" && body.Filename != "" {
		responses.Error(w, http.StatusBadRequest, "cannot change path and filename at the same time")
		return
	}

	// Get file statistics and ensure exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		responses.Error(w, http.StatusNotFound, "specified file/directory does not exist")
		return
	} else if err != nil {
		log.Printf("ERROR: failed to stat file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	// Rename file if passed
	if body.Filename != "" {
		newName := filepath.Base(filepath.Clean(body.Filename))
		directory := filepath.Dir(path)

		// Rename file
		if err := os.Rename(path, directory+"/"+newName); err != nil {
			log.Printf("ERROR: failed to rename file: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to rename file")
			return
		}
	}

	// Move file if passed
	if body.Path != "" {
		newPath := filepath.Join(filesDirectory, r.Header.Get("X-BPI-Username"), strings.Replace(body.Path, "../", "", -1))
		filename := filepath.Base(path)

		// Ensure new path exists
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			responses.Error(w, http.StatusBadRequest, "specified path does not exist")
			return
		}

		// Move file
		if err := os.Rename(path, newPath+"/"+filename); err != nil {
			log.Printf("ERROR: failed to move file to specified directory: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to move file")
			return
		}
	}

	responses.Success(w)
}

// Delete a file
func deleteFile(w http.ResponseWriter, _ *http.Request, path string) {
	// Ensure file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		responses.Error(w, http.StatusNotFound, "specified file/directory does not exist")
		return
	} else if err != nil {
		log.Printf("ERROR: failed to stat file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to stat file")
		return
	}

	// Remove all under if directory
	if info.IsDir() {
		if err := os.RemoveAll(path); err != nil {
			log.Printf("ERROR: failed to delete directory: %v\n", err)
			responses.Error(w, http.StatusInternalServerError, "failed to delete directory")
			return
		}

		responses.Success(w)
		return
	}

	// Remove file
	if err := os.Remove(path); err != nil {
		log.Printf("ERROR: failed to delete file: %v\n", err)
		responses.Error(w, http.StatusInternalServerError, "failed to remove file")
		return
	}

	responses.Success(w)
}
