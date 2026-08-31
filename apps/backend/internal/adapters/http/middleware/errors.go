package http

import (
	service "cidadon/internal/application/contract"
	"net/http"
)

type Response struct {
	Code    service.Code `json:"code"`
	Message string       `json:"message"`
	Details any          `json:"details,omitempty"`
}

var codeToStatus = map[service.Code]int{
	service.CodeInvalidInput:    http.StatusBadRequest,
	service.CodeUnauthorized:    http.StatusUnauthorized,
	service.CodeForbidden:       http.StatusForbidden,
	service.CodeNotFound:        http.StatusNotFound,
	service.CodeAlreadyExists:   http.StatusConflict,
	service.CodeConflict:        http.StatusConflict,
	service.CodeFailedPrecond:   http.StatusPreconditionFailed,
	service.CodeResourceExhaust: http.StatusTooManyRequests,
	service.CodeTimeout:         http.StatusGatewayTimeout,
	service.CodeUnavailable:     http.StatusServiceUnavailable,
	service.CodeInternal:        http.StatusInternalServerError,
}

func StatusFor(err error) int {
	if status, ok := codeToStatus[service.GetCode(err)]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func FromError(err error) (int, Response) {
	appErr, ok := service.From(err)
	if !ok {
		return http.StatusInternalServerError, Response{
			Code:    service.CodeInternal,
			Message: "internal server error",
		}
	}

	status := StatusFor(err)
	resp := Response{
		Code:    appErr.Code,
		Message: appErr.Message,
		Details: appErr.Details,
	}

	if appErr.Code == service.CodeInternal {
		resp.Details = nil
		if resp.Message == "" {
			resp.Message = "internal server error"
		}
	}

	return status, resp
}
