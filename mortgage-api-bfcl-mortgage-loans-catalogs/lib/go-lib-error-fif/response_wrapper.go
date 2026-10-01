package error_fif

import (
	"net/http"
)

func CreateResponse(code int, message string, status Status, details ...interface{}) (int, ResponseBody) {
	return code, ResponseBody{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}

func CreateBadRequest(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(400, message, BAD_REQUEST, details...)
}

func CreateUnauthorized(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusUnauthorized, message, UNAUTHORIZED, details...)
}

func CreateForbidden(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusForbidden, message, FORBIDDEN, details...)
}

func CreateNotFound(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusNotFound, message, NOT_FOUND, details...)
}

func CreateNotAcceptable(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusNotAcceptable, message, NOT_ACCEPTABLE, details...)
}

func CreateUnprocessableEntity(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusUnprocessableEntity, message, UNPROCESSABLE_ENTITY, details...)
}

func CreatePreconditionRequired(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusPreconditionRequired, message, PRECONDITION_REQUIRED, details...)
}

func CreateInternalServerError(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusInternalServerError, message, INTERNAL_ERROR, details...)
}

func CreateBadGateway(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusBadGateway, message, BAD_GATEWAY, details...)
}

func CreateGatewayTimeout(message string, details ...interface{}) (int, ResponseBody) {
	return CreateResponse(http.StatusGatewayTimeout, message, GATEWAY_TIMEOUT, details...)
}
