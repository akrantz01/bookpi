package routes

import (
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
		p := path.Join(filesDirectory, r.Header.Get("X-BPI-Username"), strings.TrimPrefix(r.RequestURI, "/api/files"))

		switch r.Method {
		case http.MethodGet:
			listFiles(w, r, p)

		case http.MethodPost:
			createFile(w, r, p)

		case http.MethodPut:
			updateFile(w, r, p)

		case http.MethodDelete:
			deleteFile(w, r, p)

		default:
			responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// List all files in a directory or a file's information
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

	responses.SuccessWithData(w, map[string]interface{}{
		"name": info.Name(),
		"size": info.Size(),
		"last_modified": info.ModTime().Unix(),
		"directory": info.IsDir(),
		"root": filepath.Clean(strings.TrimPrefix(r.RequestURI, "/api/files")) == "/",
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

	// Ensure not uploading to file
	if !info.IsDir() {
		responses.Error(w, http.StatusBadRequest, "cannot upload to file")
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
func updateFile(w http.ResponseWriter, r *http.Request, path string) {

}

// Delete a file
func deleteFile(w http.ResponseWriter, r *http.Request, path string) {

}
