package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		url := r.RequestURI
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		contentLength := r.ContentLength
		timestamp := time.Now().Format(time.RFC3339)

		green := color.New(color.FgGreen).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		orange := color.New(color.FgRed, color.Bold).SprintFunc()

		// Log thông tin yêu cầu với màu sắc
		log.Printf("[%s] %s %s - Client IP: %s - User-Agent: %s - Content-Length: %s\n",
			green(timestamp), blue(method), yellow(url), red(clientIP), orange(userAgent), orange(contentLength))

		next.ServeHTTP(w, r)
	})
}
