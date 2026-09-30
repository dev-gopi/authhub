package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/dev-gopi/authhub/internal/shared/response"
)

type HealthChecker interface {
	Health(context.Context) error
}

type HealthHandler struct {
	postgres HealthChecker
	redis    HealthChecker
	rabbitMQ HealthChecker
	vault    HealthChecker
	service  string
}

func NewHealthHandler(
	postgres HealthChecker,
	redis HealthChecker,
	rabbitMQ HealthChecker,
	vault HealthChecker,
	service string,
) *HealthHandler {

	return &HealthHandler{
		postgres: postgres,
		redis:    redis,
		rabbitMQ: rabbitMQ,
		vault:    vault,
		service:  service,
	}
}

func (h *HealthHandler) Health(
	w http.ResponseWriter,
	_ *http.Request,
) {
	response.JSON(
		w,
		http.StatusOK,
		map[string]any{
			"status":  "ok",
			"service": h.service,
		},
	)
}

func (h *HealthHandler) Ready(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		3*time.Second,
	)

	defer cancel()

	dependencies := map[string]string{
		"postgres": "ok",
		"redis":    "ok",
		"rabbitmq": "ok",
		"vault":    "ok",
	}

	ready := true

	if err := h.postgres.Health(ctx); err != nil {
		dependencies["postgres"] = "unavailable"
		ready = false
	}

	if err := h.redis.Health(ctx); err != nil {
		dependencies["redis"] = "unavailable"
		ready = false
	}

	if err := h.rabbitMQ.Health(ctx); err != nil {
		dependencies["rabbitmq"] = "unavailable"
		ready = false
	}

	if err := h.vault.Health(ctx); err != nil {
		dependencies["vault"] = "unavailable"
		ready = false
	}

	status := http.StatusOK
	state := "ready"

	if !ready {
		status = http.StatusServiceUnavailable
		state = "not_ready"
	}

	response.JSON(
		w,
		status,
		map[string]any{
			"status":       state,
			"service":      h.service,
			"dependencies": dependencies,
		},
	)
}
