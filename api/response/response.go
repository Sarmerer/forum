package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type JSONError struct {
	Error  string
	Status int
}

func Error(w http.ResponseWriter, errorCode int, err error) {
	if err != nil {
		w.WriteHeader(errorCode)
		JSON(w, err.Error(), errorCode)
	}
}

func JSON(w http.ResponseWriter, err string, errCode int) {
	b, marshalErr := json.Marshal(JSONError{Error: fmt.Sprint(err), Status: errCode})
	if marshalErr != nil {
		b = []byte("internal server error.")
	}
	w.Write(b)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, errors.New("internal server error"))
}

func BadRequest(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, errors.New("bad request"))
}
