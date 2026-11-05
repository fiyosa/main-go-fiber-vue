package exceptions

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/response"
)

// Handler is the global Fiber error handler. It is the Go counterpart of the
// Laravel ApiExceptionHandler trait: map the error, then write one JSON
// envelope with a translated message.
func Handler(debug bool) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		ex := FromError(err)

		if ex.Status >= fiber.StatusInternalServerError {
			slog.Error("request failed", "code", ex.Code, "status", ex.Status, "error", ex.Err)
		}

		if ex.Fields != nil {
			return response.ValidationError(c, ex.Fields)
		}
		if debug && ex.Err != nil && ex.Status >= fiber.StatusInternalServerError {
			return response.ErrorDebug(c, ex.MessageKey, ex.Status, ex.Err.Error())
		}
		return response.Error(c, ex.MessageKey, nil, ex.Status)
	}
}

// FromError resolves any error into an Exception. Known exceptions are kept,
// Fiber transport errors are mapped by status, and everything else becomes a
// 500 with a translated generic message.
func FromError(err error) *Exception {
	if ex, ok := As(err); ok {
		return ex
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fromHTTPStatus(fiberErr.Code, fiberErr)
	}

	return Internal(err)
}

func fromHTTPStatus(status int, cause error) *Exception {
	switch status {
	case fiber.StatusBadRequest:
		return BadRequest(cause)
	case fiber.StatusUnauthorized:
		return Unauthorized("auth.unauthenticated", cause)
	case fiber.StatusForbidden:
		return Forbidden(cause)
	case fiber.StatusNotFound:
		return NotFound(cause)
	case fiber.StatusMethodNotAllowed:
		return MethodNotAllowed(cause)
	case fiber.StatusUnsupportedMediaType:
		return New("common.bad_request", status, "common.bad_request", cause)
	case fiber.StatusServiceUnavailable:
		return ServiceUnavailable(cause)
	}
	return New("http.error", status, "common.try_again_later", cause)
}
