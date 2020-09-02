package helpers

import (
	"errors"
	"net/http"
	"strconv"
)

func ParseID(r *http.Request) (res uint64, err error) {
	if res, err = strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64); err != nil {
		return 0, errors.New("invalid id")
	}
	return res, nil
}
