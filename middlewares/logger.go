package middlewares

import (
	"log"
	"net/http"
	"time"
)

// Logger middleware untuk log request masuk dan durasi eksekusi
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("[REQUEST] %s %s dari %s", r.Method, r.RequestURI, r.RemoteAddr)

		// Jalankan handler
		next(w, r)

		duration := time.Since(start)
		log.Printf("[DONE]    %s %s selesai dalam %v", r.Method, r.RequestURI, duration)
	}
}
