package handler

import (
	"net/http"
	"web-service/pkg/utils"

	"github.com/gorilla/mux"
)

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	utils.ResponseJson(
		w,
		http.StatusNotFound,
		utils.Response{
			StatusCode: http.StatusNotFound,
			Message:    "404 Not Found",
			Data:       nil,
		},
	)
}

func NotFoundHandler(r *mux.Router) {
	r.NotFoundHandler = http.HandlerFunc(notFoundResponse)

}
