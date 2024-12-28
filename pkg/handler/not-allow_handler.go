package handler

import (
	"net/http"
	"web-service/pkg/utils"

	"github.com/gorilla/mux"
)

func notAllowResponse(w http.ResponseWriter, r *http.Request) {
	utils.ResponseJson(
		w,
		http.StatusMethodNotAllowed,
		utils.Response{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "404 Not Found",
			Data:       nil,
		},
	)
}

func NotAllowHandler(r *mux.Router) {
	r.MethodNotAllowedHandler = http.HandlerFunc(notAllowResponse)
}
