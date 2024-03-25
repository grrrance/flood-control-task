package utils

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"task/pkg/logger"
)

func GetRequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

func GetRequestIP(r *http.Request) string {
	return r.RemoteAddr
}

func LogResponseError(r *http.Request, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(r.Context()),
		GetRequestIP(r),
		err.Error(),
	)
}

func HandleResponse(w http.ResponseWriter, code int, resp interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
