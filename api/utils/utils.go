package utils

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ParseID(r *http.Request) (res uint64, err error) {
	if res, err = strconv.ParseUint(r.URL.Query().Get("ID"), 10, 64); err != nil {
		return 0, errors.New("invalid id")
	}
	return res, nil
}

// LoadEnv sets environment variables required to run the API
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	return parse(file)
}

func parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		envVar := strings.Split(scanner.Text(), "=")
		os.Setenv(envVar[0], envVar[1])
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if key := os.Getenv("API_KEY"); key == "" {
		return errors.New("could not find API_KEY environment variable")
	}
	return nil
}
