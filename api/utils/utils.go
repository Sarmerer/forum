package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sarmerer/forum/api/config"

	"github.com/sarmerer/forum/api/models"
)

func ParseID(r *http.Request) (res int64, err error) {
	if res, err = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64); err != nil {
		return 0, errors.New("invalid id")
	}
	return res, nil
}

func SetupEnv() error {
	var (
		file      *os.File
		lineCount int = 1
		err       error
	)
	if file, err = os.Open("./.env"); err != nil {
		return errors.New("could not find .env file, skipping...")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		env := strings.Split(scanner.Text(), "=")
		if len(env) < 2 {
			return fmt.Errorf("invalid env variable on line %d", lineCount)
		}
		lineCount++
		if os.Getenv(env[0]) != "" {
			break
		}
		os.Setenv(env[0], env[1])
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	return nil
}

func GetUserFromCtx(r *http.Request) models.UserCtx {
	if r.Context().Value(config.UserCtxVarName) != nil {
		return r.Context().Value(config.UserCtxVarName).(models.UserCtx)
	}
	return models.UserCtx{ID: -1, Role: -1}
}

func CurrentUnixTime() int64 {
	return time.Now().Unix()
}

func ParseFlags(args []string) []string {
	var res []string
	for _, arg := range args {
		if _, ok := config.Flags[arg]; ok {
			*config.Flags[arg].State = true
			res = append(res, config.Flags[arg].Message)
		}
	}
	return res
}

func Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

func FormatRequest(r *http.Request) string {
	var request []string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	return strings.Join(request, "\n")
}
