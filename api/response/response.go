package response

import (
	"encoding/json"
	"errors"
	"net/http"
)

type response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func Error(w http.ResponseWriter, errorCode int, err error) {
	if err != nil {
		w.WriteHeader(errorCode)
		JSON(w, "error", errorCode, err.Error(), nil)
	}
}

func JSON(w http.ResponseWriter, status string, code int, message, data interface{}) {
	b, marshalErr := json.Marshal(response{status, code, message, data})
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b = []byte(`{"Error":"internal server error", "Status":500}`)
	}
	w.Write(b)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, errors.New("internal server error"))
}

func BadRequest(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, errors.New("bad request"))
}
