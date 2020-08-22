package utils

import (
	"bufio"
	"io"
	"os"
	"strings"
)

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
	return nil
}
