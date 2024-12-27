package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"web-service/pkg/data"
	"web-service/pkg/utils"

	"github.com/gorilla/mux"
)

func getProducts(w http.ResponseWriter, r *http.Request) utils.Response {

	return utils.SuccessResponse("Get all products successfully", data.ListProduct)
}

func getProductById(w http.ResponseWriter, r *http.Request) utils.Response {
	vars := mux.Vars(r)
	id, err := utils.GetId(vars["id"])

	if err != nil {
		return utils.BadRequestError("Invalid id", nil)
	}

	for _, product := range data.ListProduct {
		if product.ID == id {
			return utils.SuccessResponse(fmt.Sprintf("Get product with id %d successfully", id), product)
		}
	}

	return utils.NotFoundError(fmt.Sprintf("Product with id %d not found", id), nil)
}

func createProduct(w http.ResponseWriter, r *http.Request) utils.Response {
	defer r.Body.Close()

	var product data.ProductData

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	if err := decode.Decode(&product); err != nil {
		errorMessage := utils.JSONDecodeError(err)

		return utils.BadRequestError(errorMessage, nil)
	}

	data.ListProduct = append(data.ListProduct, product)

	return utils.CreatedResponse("Create product successfully", product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) utils.Response {
	vars := mux.Vars(r)
	id, err := utils.GetId(vars["id"])

	if err != nil {
		utils.BadRequestError("Invalid id", nil)
	}

	var result []data.ProductData
	var found bool = false

	for _, product := range data.ListProduct {
		if product.ID != id {
			result = append(result, product)
		} else {
			found = true
		}
	}

	if !found {
		return utils.NotFoundError(fmt.Sprintf("Product with id %d not found", id), nil)
	}

	data.ListProduct = result

	return utils.SuccessResponse(fmt.Sprintf("Delete product with id %d successfully", id), nil)
}

func ProductRoutes(r *mux.Router) {
	productRouter := r.PathPrefix("/products").Subrouter().StrictSlash(true)

	productRouter.HandleFunc("/", utils.WrapHandler(getProducts)).Methods(http.MethodGet).Name("getProducts")
	productRouter.HandleFunc("/{id:[0-9]+}", utils.WrapHandler(getProductById)).Methods(http.MethodGet).Name("getProductById")
	productRouter.HandleFunc("/create", utils.WrapHandler(createProduct)).Methods(http.MethodPost).Name("createProduct")
	productRouter.HandleFunc("/{id:[0-9]+}", utils.WrapHandler(deleteProduct)).Methods(http.MethodDelete).Name("deleteProduct")
}
