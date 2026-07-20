package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler — глобальный hook ogen. Подключается через orderv1.WithErrorHandler.
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
	case errors.As(err, &decodeParams), errors.As(err, &decodeRequest):
		return http.StatusBadRequest, err.Error()

	// 404 Not Found
	case errors.Is(err, errs.ErrOrderNotFound),
		errors.Is(err, errs.ErrOrderItemNotFound),
		errors.Is(err, errs.ErrPartNotFound):
		return http.StatusNotFound, err.Error()

	// 409 Conflict
	case errors.Is(err, errs.ErrOrderAlreadyPaid),
		errors.Is(err, errs.ErrOrderCancelled),
		errors.Is(err, errs.ErrOutOfStock),
		errors.Is(err, errs.ErrOrderStatusIncorrect):
		return http.StatusConflict, err.Error()

	// 400 Bad Request
	case errors.Is(err, errs.ErrInvalidUUID),
		errors.Is(err, errs.ErrInvalidPaymentMethod),
		errors.Is(err, errs.ErrMissingRequiredParts):
		return http.StatusBadRequest, err.Error()

	// 500 Internal Server Error
	default:
		return http.StatusInternalServerError, "внутренняя ошибка"
	}
}
