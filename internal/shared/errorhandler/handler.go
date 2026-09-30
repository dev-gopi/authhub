package errorhandler

import (
	"errors"
	"net/http"

	"github.com/dev-gopi/authhub/internal/shared/response"
)

func Handle(w http.ResponseWriter, err error) {
	var appErr *AppError

	if errors.As(err, &appErr) {
		response.Error(
			w,
			appErr.HTTPStatus,
			appErr.Code,
			appErr.Message,
		)
		return
	}

	response.Error(
		w,
		http.StatusInternalServerError,
		CodeInternalError,
		"internal server error",
	)
}
