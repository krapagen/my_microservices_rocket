package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler — глобальный hook ogen. Подключается через iamv1.WithErrorHandler.
func ErrorHandler(ctx context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code, message := mapError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if encErr := json.NewEncoder(w).Encode(errorResponse{Code: code, Message: message}); encErr != nil {
		slog.ErrorContext(ctx, "ошибка кодирования ответа", "error", encErr)
	}
}

func mapError(err error) (int, string) {
	// Ошибки декодирования/валидации запроса от ogen — всегда 400.
	var decodeParams *ogenerrors.DecodeParamsError
	var decodeRequest *ogenerrors.DecodeRequestError

	switch {
	// 400 Bad Request — декодирование и доменная валидация входа
	case errors.As(err, &decodeParams), errors.As(err, &decodeRequest):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, errs.ErrInvalidLogin),
		errors.Is(err, errs.ErrWeakPassword),
		errors.Is(err, errs.ErrEmptyCredential),
		errors.Is(err, errs.ErrEmptySessionID),
		errors.Is(err, errs.ErrInvalidUUID):
		return http.StatusBadRequest, err.Error()

	// 401 Unauthorized — неверные креды или недействительная сессия
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrSessionNotFound):
		return http.StatusUnauthorized, err.Error()

	// 404 Not Found
	case errors.Is(err, errs.ErrUserNotFound):
		return http.StatusNotFound, err.Error()

	// 409 Conflict
	case errors.Is(err, errs.ErrUserAlreadyExists):
		return http.StatusConflict, err.Error()

	// 500 Internal Server Error
	default:
		return http.StatusInternalServerError, "внутренняя ошибка"
	}
}
