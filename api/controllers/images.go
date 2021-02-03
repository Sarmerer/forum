package controllers

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sarmerer/forum/api/config"
	"github.com/sarmerer/forum/api/response"
	uuid "github.com/satori/go.uuid"
)

func ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("image")
	if filename == "" {
		response.Error(w, http.StatusBadRequest, errors.New("no image name provided"))
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	http.ServeFile(w, r, fmt.Sprintf("./database/images/%s", filename))
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	var (
		maxImageSize int64 = config.MaxImageUploadSize
		image        multipart.File
		handler      *multipart.FileHeader
		file         *os.File
		fileName     string
		err          error
	)
	if r.ContentLength > maxImageSize {
		response.Error(w, http.StatusExpectationFailed, errors.New("image is too heavy"))
		return
	}
	if image, handler, err = r.FormFile("image"); err != nil {
		response.Error(w, http.StatusBadRequest, err)
		return
	}
	defer image.Close()
	defer file.Close()

	fileNameSplit := strings.Split(handler.Filename, ".")
	if len(fileNameSplit) < 1 {
		response.Error(w, http.StatusBadRequest, errors.New("invalid image extension"))
		return
	}
	extension := fileNameSplit[len(fileNameSplit)-1]
	fileName = fmt.Sprint(uuid.NewV4())

	if file, err = os.Create(fmt.Sprintf("./database/images/%s.%s", fileName, extension)); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if _, err = io.Copy(file, image); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "file has been uploaded", fmt.Sprintf("%s/images?image=%s.%s", config.APIURL, fileName, extension))
}
