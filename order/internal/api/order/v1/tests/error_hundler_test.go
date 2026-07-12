package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
)

type failWriter struct {
	header http.Header
	code   int
}

func (w *failWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *failWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func (w *failWriter) WriteHeader(code int) {
	w.code = code
}

func (s *ServiceSuite) TestErrorHandler() {
	tests := []struct {
		name       string
		inputErr   error
		wantCode   int
		wantInBody string
	}{
		{
			name:       "decode params error",
			inputErr:   &ogenerrors.DecodeParamsError{OperationContext: ogenerrors.OperationContext{Name: "test"}, Err: errors.New("bad params")},
			wantCode:   http.StatusBadRequest,
			wantInBody: "bad params",
		},
		{
			name:       "decode request error",
			inputErr:   &ogenerrors.DecodeRequestError{OperationContext: ogenerrors.OperationContext{Name: "test"}, Err: errors.New("bad request")},
			wantCode:   http.StatusBadRequest,
			wantInBody: "bad request",
		},
		{
			name:     "order not found",
			inputErr: errs.ErrOrderNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "part not found",
			inputErr: errs.ErrPartNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "order already paid",
			inputErr: errs.ErrOrderAlreadyPaid,
			wantCode: http.StatusConflict,
		},
		{
			name:     "order cancelled",
			inputErr: errs.ErrOrderCancelled,
			wantCode: http.StatusConflict,
		},
		{
			name:     "out of stock",
			inputErr: errs.ErrOutOfStock,
			wantCode: http.StatusConflict,
		},
		{
			name:     "order status incorrect",
			inputErr: errs.ErrOrderStatusIncorrect,
			wantCode: http.StatusConflict,
		},
		{
			name:     "invalid uuid",
			inputErr: errs.ErrInvalidUUID,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid payment method",
			inputErr: errs.ErrInvalidPaymentMethod,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "missing required parts",
			inputErr: errs.ErrMissingRequiredParts,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "generic error",
			inputErr: errors.New("boom"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			v1.ErrorHandler(s.ctx, rec, req, tt.inputErr)

			s.Equal(tt.wantCode, rec.Code)
			s.Equal("application/json", rec.Header().Get("Content-Type"))
			s.Contains(rec.Body.String(), tt.wantInBody)
		})
	}
}

func (s *ServiceSuite) TestErrorHandler_EncodeError() {
	fw := &failWriter{}
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	v1.ErrorHandler(s.ctx, fw, req, errs.ErrOrderNotFound)

	s.Equal(http.StatusNotFound, fw.code)
}
