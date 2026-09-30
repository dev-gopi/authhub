package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool `json:"success"`

	Data any `json:"data,omitempty"`

	Error any `json:"error,omitempty"`
}

func JSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func Success(
	w http.ResponseWriter,
	status int,
	data any,
) {
	JSON(w, status, Response{
		Success: true,
		Data:    data,
	})
}

func Error(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	JSON(w, status, Response{
		Success: false,
		Error: map[string]any{
			"code":    code,
			"message": message,
		},
	})
}
