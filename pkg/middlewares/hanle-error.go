package middlewares

import (
	"log"
	"net/http"
	"web-service/pkg/utils"
)

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				response := utils.Response{
					StatusCode: http.StatusInternalServerError,
					Message:    "Internal Server Error",
					Data:       nil,
				}

				utils.ResponseJson(w, http.StatusInternalServerError, response)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
