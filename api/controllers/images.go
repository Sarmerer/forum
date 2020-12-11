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

	"github.com/sarmerer/forum/api/response"
	uuid "github.com/satori/go.uuid"
)

func ServeImage(w http.ResponseWriter, r *http.Request) {
	filename := "logo.png"
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
	http.ServeFile(w, r, "./database/images/"+filename)
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	var (
		maxImageSize int64 = 5 * 1024 * 1024 // 5mb
		image        multipart.File
		handler      *multipart.FileHeader
		file         *os.File
		fileName     string
		err          error
	)
	if r.ContentLength > maxImageSize {
		response.Error(w, http.StatusBadRequest, errors.New("image is too heavy"))
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

	if file, err = os.Create("./database/images/" + fileName + extension); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}

	if _, err = io.Copy(file, image); err != nil {
		response.Error(w, http.StatusInternalServerError, err)
		return
	}
	response.Success(w, "file has been uploaded", fileName)
}
