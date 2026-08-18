package errors

import (
	"net/http"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

var codeToStatus = map[Code]int{
	CodeInvalidInput:    http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,
	CodeAlreadyExists:   http.StatusConflict,
	CodeConflict:        http.StatusConflict,
	CodeFailedPrecond:   http.StatusPreconditionFailed,
	CodeResourceExhaust: http.StatusTooManyRequests,
	CodeTimeout:         http.StatusGatewayTimeout,
	CodeUnavailable:     http.StatusServiceUnavailable,
	CodeInternal:        http.StatusInternalServerError,
}

func StatusFor(err error) int {
	if status, ok := codeToStatus[GetCode(err)]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func FromError(err error) (int, Response) {
	appErr, ok := From(err)
	if !ok {
		return http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}

	status := StatusFor(err)
	resp := Response{
		Code:    status,
		Message: appErr.Message,
		Details: appErr.Details,
	}

	if appErr.Code == CodeInternal {
		resp.Details = nil
		if resp.Message == "" {
			resp.Message = "internal server error"
		}
	}

	return status, resp
}
