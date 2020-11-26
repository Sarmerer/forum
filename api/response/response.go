package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

// JSON builds and writes an API response. Arguments:
// 	responseStatus  - defines the status of a response, can be success or error
// 	httpStatus      - http code of a response
//	message 		- provides additional information about a response
// data				- body of a response
func JSON(w http.ResponseWriter, responseStatus string, httpStatus int, message, data interface{}) {
	////fmt.Println(responseStatus, httpStatus, message, data)
	b, marshalErr := json.Marshal(response{responseStatus, httpStatus, message, data})
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		b = []byte(`{"status":"error", "code":500,"message":"failed to marshal JSON in response.JSON()","data":null}`)
	}
	w.Write(b)
}

// Error acts as a template for a JSON function.
// It automatically sets status to "error"
func Error(w http.ResponseWriter, httpStatus int, err error) {
	fmt.Println(err)
	w.WriteHeader(httpStatus)
	JSON(w, "error", httpStatus, err.Error(), nil)
}

// Success acts as a template for a JSON function.
// It automatically sets response status to "success" and http status to 200
func Success(w http.ResponseWriter, message, data interface{}) {
	w.WriteHeader(http.StatusOK)
	JSON(w, "success", http.StatusOK, message, data)
}
