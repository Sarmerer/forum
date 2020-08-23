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

func SetHeaders(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		//CORS handler
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next(w, r)
	}
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, 200, time.Now()}
		next(&rec, r)
		log.Printf("|%s|\t%10s|\t%s |%s %s", utils.PaintStatus(rec.status), time.Since(rec.elapsed), r.Host, utils.PaintMethod(r.Method), r.URL.Path)
	}
}
