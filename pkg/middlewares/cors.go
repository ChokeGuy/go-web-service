package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"web-service/config"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", config.Env.CORSOrigins[0])
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Env.CORSMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Env.CORSHeaders, ","))
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.Env.CORSMaxAge))

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
