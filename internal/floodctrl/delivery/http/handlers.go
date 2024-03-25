package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"task/internal/floodctrl"
	"task/pkg/logger"
	"task/pkg/utils"
)

type FloodResponse struct {
	IsTriggered bool `json:"is_triggered"`
}

type floodHandlers struct {
	floodUC floodctrl.FloodControl
	logger  logger.Logger
}

func NewFloodHandlers(floodUC floodctrl.FloodControl, log logger.Logger) floodctrl.Handlers {
	return &floodHandlers{floodUC: floodUC, logger: log}
}

func (f *floodHandlers) TriggerUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resp := &FloodResponse{IsTriggered: false}
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			utils.LogResponseError(r, f.logger, err)
			utils.HandleResponse(w, http.StatusInternalServerError, resp)
			return
		}

		isTriggered, err := f.floodUC.Check(ctx, id)
		if err != nil {
			utils.LogResponseError(r, f.logger, err)
			utils.HandleResponse(w, http.StatusInternalServerError, resp)
			return
		}
		resp.IsTriggered = isTriggered
		if resp.IsTriggered {
			utils.HandleResponse(w, http.StatusOK, resp)
		} else {
			utils.HandleResponse(w, http.StatusTooManyRequests, resp)
		}
	}
}
