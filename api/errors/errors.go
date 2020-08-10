package errors

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	URL    string
	Status int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("There was an error, while serving %s. Error code: %d\n", e.URL, e.Status)
}

func HTTPErrorsHandler(errorCode int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(errorCode)
	w.Write([]byte(fmt.Sprint(errorCode)))
}
