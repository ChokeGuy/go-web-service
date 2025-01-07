package utils

import "net/http"

func (e *Response) Error() string {
	return e.Message
}

func BadRequestError(message string, data interface{}) Response {
	return Response{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Data:       data,
	}
}

func NotFoundError(message string, data interface{}) Response {
	return Response{
		StatusCode: http.StatusNotFound,
		Message:    message,
		Data:       data,
	}
}

func ValidationError(message string, data interface{}) Response {
	return Response{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
		Data:       data,
	}
}

func InternalServerError(message string) Response {
	return Response{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Data:       nil,
	}
}

func UnauthorizedError(message string) Response {
	return Response{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Data:       nil,
	}
}
