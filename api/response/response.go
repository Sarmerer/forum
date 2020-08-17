package response

import (
	"encoding/json"
	"errors"
	"forum/config"
	"net/http"
)

type response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

//JSON function takes
func JSON(w http.ResponseWriter, status string, code int, message, data interface{}) {
	b, marshalErr := json.Marshal(response{status, code, message, data})
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b = []byte(`{"status":"error", "code":500,"message":"internal server error","data":null}`)
	}
	w.Write(b)
}

func Error(w http.ResponseWriter, httpStatus int, err error) {
	w.WriteHeader(httpStatus)
	JSON(w, config.StatusError, httpStatus, err.Error(), nil)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, errors.New("internal server error"))
}

func BadRequest(w http.ResponseWriter) {
	Error(w, http.StatusBadRequest, errors.New("bad request"))
}
