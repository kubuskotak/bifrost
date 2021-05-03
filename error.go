package bifrost

import (
	"context"
	"net/http"
)

type ctxError struct {
	Name string
}

func (r *ctxError) String() string {
	return "context value " + r.Name
}

var CtxError = ctxError{Name: "context error"}

// ErrBadRequest error http StatusBadRequest
func ErrBadRequest(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusBadRequest))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadRequest)
	return err
}

// ErrUnauthorized error http StatusUnauthorized
func ErrUnauthorized(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusUnauthorized))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnauthorized)
	return err
}

// ErrPaymentRequired error http StatusPaymentRequired
func ErrPaymentRequired(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusPaymentRequired))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusPaymentRequired)
	return err
}

// ErrForbidden error http StatusForbidden
func ErrForbidden(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusForbidden))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusForbidden)
	return err
}

// ErrMethodNotAllowed error http StatusMethodNotAllowed
func ErrMethodNotAllowed(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusMethodNotAllowed))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusMethodNotAllowed)
	return err
}

// ErrNotAcceptable error http StatusNotAcceptable
func ErrNotAcceptable(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusNotAcceptable))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotAcceptable)
	return err
}

// ErrProxyAuthRequired error http StatusProxyAuthRequired
func ErrProxyAuthRequired(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusProxyAuthRequired))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusProxyAuthRequired)
	return err
}

// ErrRequestTimeout error http StatusRequestTimeout
func ErrRequestTimeout(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusRequestTimeout))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusRequestTimeout)
	return err
}

// ErrUnsupportedMediaType error http StatusUnsupportedMediaType
func ErrUnsupportedMediaType(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusUnsupportedMediaType))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnsupportedMediaType)
	return err
}

// ErrUnprocessableEntity error http StatusUnprocessableEntity
func ErrUnprocessableEntity(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusUnprocessableEntity))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnprocessableEntity)
	return err
}

// ErrInternalServerError error http StatusInternalServerError
func ErrInternalServerError(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusInternalServerError))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	return err
}

// ErrBadGateway error http StatusBadGateway
func ErrBadGateway(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusBadGateway))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadGateway)
	return err
}

// ErrServiceUnavailable error http StatusServiceUnavailable
func ErrServiceUnavailable(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusServiceUnavailable))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusServiceUnavailable)
	return err
}

// ErrGatewayTimeout error http StatusGatewayTimeout
func ErrGatewayTimeout(w http.ResponseWriter, r *http.Request, err error) error {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxError, http.StatusGatewayTimeout))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusGatewayTimeout)
	return err
}
