package handler

import (
	"net/http"
	"web-service/pkg/utils"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) utils.Response {
	return utils.SuccessResponse("Welcome to the home page", nil)

}

func HomeRoutes(r *mux.Router) {
	r.HandleFunc("/", utils.WrapHandler(homeHandler)).Methods("GET")
}
