package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/ogen-go/ogen/ogenerrors"

	iamv1api "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/iam/v1"
	errs "github.com/krapagen/my_microservices_rocket/iam/internal/errors"
)

func (s *IamAPISuite) callErrorHandler(err error) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	iamv1api.ErrorHandler(s.ctx, rec, req, err)
	return rec
}

func (s *IamAPISuite) TestErrorHandler_BadRequest() {
	for _, err := range []error{
		errs.ErrInvalidLogin,
		errs.ErrWeakPassword,
		errs.ErrEmptyCredential,
		errs.ErrEmptySessionID,
		errs.ErrInvalidUUID,
	} {
		rec := s.callErrorHandler(err)

		s.Equal(http.StatusBadRequest, rec.Code, "ошибка: %v", err)
		s.Equal("application/json", rec.Header().Get("Content-Type"))

		var body struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		s.NoError(json.NewDecoder(rec.Body).Decode(&body))
		s.Equal(http.StatusBadRequest, body.Code)
		s.NotEmpty(body.Message)
	}
}

func (s *IamAPISuite) TestErrorHandler_DecodeErrors() {
	// Ошибки декодирования от ogen — всегда 400
	decodeErrs := []error{
		&ogenerrors.DecodeParamsError{Err: errors.New("bad param")},
		&ogenerrors.DecodeRequestError{Err: errors.New("bad body")},
	}

	for _, err := range decodeErrs {
		rec := s.callErrorHandler(err)
		s.Equal(http.StatusBadRequest, rec.Code, "ошибка: %v", err)
	}
}

func (s *IamAPISuite) TestErrorHandler_Unauthorized() {
	for _, err := range []error{
		errs.ErrInvalidCredentials,
		errs.ErrSessionNotFound,
	} {
		rec := s.callErrorHandler(err)
		s.Equal(http.StatusUnauthorized, rec.Code, "ошибка: %v", err)
	}
}

func (s *IamAPISuite) TestErrorHandler_NotFound() {
	rec := s.callErrorHandler(errs.ErrUserNotFound)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *IamAPISuite) TestErrorHandler_Conflict() {
	rec := s.callErrorHandler(errs.ErrUserAlreadyExists)
	s.Equal(http.StatusConflict, rec.Code)
}

func (s *IamAPISuite) TestErrorHandler_Internal() {
	rec := s.callErrorHandler(errors.New("unexpected"))

	s.Equal(http.StatusInternalServerError, rec.Code)

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	s.NoError(json.NewDecoder(rec.Body).Decode(&body))
	s.Equal("внутренняя ошибка", body.Message)
}

// failingResponseWriter — ResponseWriter, падающий на записи тела.
// Используется для покрытия ветки ошибки кодирования в ErrorHandler
type failingResponseWriter struct {
	header http.Header
}

func (w *failingResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *failingResponseWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func (w *failingResponseWriter) WriteHeader(int) {}

// TestErrorHandler_EncodeError — ошибка записи ответа не должна паниковать:
// ErrorHandler логирует её и завершается
func (s *IamAPISuite) TestErrorHandler_EncodeError() {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	s.NotPanics(func() {
		iamv1api.ErrorHandler(s.ctx, &failingResponseWriter{}, req, errors.New("boom"))
	})
}
