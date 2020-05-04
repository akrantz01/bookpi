package assets

import (
	"fmt"
	"net/http"
	"os"
)

var StaticServer = http.FileServer(new(staticServer))
var notFound http.File

func init() {
	// Load into memory
	var err error
	notFound, err = FS.OpenFile(CTX, "404.html", os.O_RDONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open 404.html: %v", err))
	}
}

type staticServer struct {}

func (s *staticServer) Open(path string) (http.File, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return notFound, nil
	} else if err != nil {
		return nil, err
	}
	return f, nil
}
