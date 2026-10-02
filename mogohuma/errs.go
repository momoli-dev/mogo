package mogohuma

import (
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

// Err400 returns 400 Bad Request with public error details.
func Err400(err ...error) error {
	return huma.Error400BadRequest("bad request", err...)
}

// Err401 returns 401 Unauthorized, logging err at debug level without exposing it.
func Err401(err error) error {
	slog.Debug("unauthorized", "error", err)
	return huma.Error401Unauthorized("unauthorized")
}

// Err403 returns 403 Forbidden, logging err at debug level without exposing it.
func Err403(err error) error {
	slog.Debug("forbidden", "error", err)
	return huma.Error403Forbidden("forbidden")
}

// Err404 returns 404 Not Found with public error details.
func Err404(err error) error {
	return huma.Error404NotFound("not found", err)
}

// Err500 returns 500 Internal Server Error, logging err at error level without exposing it.
func Err500(err error) error {
	slog.Error("internal server error", "error", err)
	return huma.Error500InternalServerError("internal server error")
}
