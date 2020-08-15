package response

import (
	"encoding/json"
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
		b = []byte("Internal server error.")
	}
	w.Write(b)
}
