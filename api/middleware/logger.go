package middleware

import (
	"forum/api/utils"
	"log"
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
		log.Printf("|%s|\t%10s|\t%s |%s %s", utils.PaintStatus(rec.status), time.Since(rec.elapsed), r.Host, utils.PaintMethod(r.Method), r.URL.Path)
	}
}
