package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPError struct {
	URL    string
	Status int
}

type JSONError struct {
	Error  string
	Status int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("There was an error, while serving %s. Error code: %d\n", e.URL, e.Status)
}

func HTTPErrors(errorCode int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(errorCode)
	w.Write([]byte(fmt.Sprint(errorCode)))
}

func JSONErrors(err string, errCode int, w http.ResponseWriter, r *http.Request) {
	b, marshalErr := json.Marshal(JSONError{Error: fmt.Sprint(err), Status: errCode})
	if marshalErr != nil {
		b = []byte("Internal server error.")
	}
	w.Write(b)
}
