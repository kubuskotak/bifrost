package bifrost

import "net/http"

// ErrBadRequest error http StatusBadRequest
func ErrBadRequest(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadRequest)
	return err
}

// ErrUnauthorized error http StatusUnauthorized
func ErrUnauthorized(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnauthorized)
	return err
}

// ErrPaymentRequired error http StatusPaymentRequired
func ErrPaymentRequired(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusPaymentRequired)
	return err
}

// ErrForbidden error http StatusForbidden
func ErrForbidden(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusForbidden)
	return err
}

// ErrMethodNotAllowed error http StatusMethodNotAllowed
func ErrMethodNotAllowed(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusMethodNotAllowed)
	return err
}

// ErrNotAcceptable error http StatusNotAcceptable
func ErrNotAcceptable(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotAcceptable)
	return err
}

// ErrProxyAuthRequired error http StatusProxyAuthRequired
func ErrProxyAuthRequired(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusProxyAuthRequired)
	return err
}

// ErrRequestTimeout error http StatusRequestTimeout
func ErrRequestTimeout(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusRequestTimeout)
	return err
}

// ErrUnsupportedMediaType error http StatusUnsupportedMediaType
func ErrUnsupportedMediaType(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnsupportedMediaType)
	return err
}

// ErrUnprocessableEntity error http StatusUnprocessableEntity
func ErrUnprocessableEntity(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnprocessableEntity)
	return err
}

// ErrInternalServerError error http StatusInternalServerError
func ErrInternalServerError(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	return err
}

// ErrBadGateway error http StatusBadGateway
func ErrBadGateway(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusBadGateway)
	return err
}

// ErrServiceUnavailable error http StatusServiceUnavailable
func ErrServiceUnavailable(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusServiceUnavailable)
	return err
}

// ErrGatewayTimeout error http StatusGatewayTimeout
func ErrGatewayTimeout(w http.ResponseWriter, err error) error {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusGatewayTimeout)
	return err
}
