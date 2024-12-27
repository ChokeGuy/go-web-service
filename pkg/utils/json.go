package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) Response

func ResponseJson(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	json.NewEncoder(writer).Encode(object)
}

func JSONDecodeError(err error) string {
	if syntaxErr, ok := err.(*json.SyntaxError); ok {
		return fmt.Sprintf("Syntax error at byte offset %d", syntaxErr.Offset)
	}
	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		return fmt.Sprintf("Type error: field %q, expected %s but got %s",
			typeErr.Field, typeErr.Type, typeErr.Value)
	}
	if strings.Contains(err.Error(), "unknown field") {
		return err.Error()
	}
	return "Invalid JSON format"
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	}
}

func CreatedResponse(message string, data interface{}) Response {
	return Response{
		StatusCode: http.StatusCreated,
		Message:    message,
		Data:       data,
	}
}

func WrapHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := handler(w, r)
		ResponseJson(w, response.StatusCode, response)
	}
}
