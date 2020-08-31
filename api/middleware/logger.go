package middleware

import (
	"fmt"
	"forum/api/logger"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status  int
	elapsed time.Time
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, 200, time.Now()}
		next(&rec, r)
		logger.HTTPLogs(logger.PaintStatus(rec.status), fmt.Sprint(time.Since(rec.elapsed)), r.Host, logger.PaintMethod(r.Method), r.URL.Path)
	}
}
