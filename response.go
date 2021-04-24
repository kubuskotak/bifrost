package bifrost

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type ctxKeyVersion struct {
	Name string
}

func (r *ctxKeyVersion) String() string {
	return "context value " + r.Name
}

type ctxKeyResponse struct {
	Name string
}

func (r *ctxKeyResponse) String() string {
	return "context value " + r.Name
}

var (
	CtxResponse = ctxKeyResponse{Name: "context Respond"}
	CtxVersion  = ctxKeyVersion{Name: "context version"}
)

type Meta struct {
	Code    string `json:"code,omitempty"`
	Type    string `json:"error_type,omitempty"`
	Message string `json:"error_message,omitempty"`
}

type Version struct {
	Label  string `json:"label,omitempty"`
	Number string `json:"number,omitempty"`
}

type Response struct {
	Version    interface{} `json:"version,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func NewResponse(r *http.Request) *Response {
	null := make(map[string]interface{})
	resp := &Response{
		Version: Version{
			Label:  "v1",
			Number: "0.1.0",
		},
		Meta:       null,
		Data:       null,
		Pagination: null,
	}
	if ver, ok := r.Context().Value(CtxVersion).(Version); ok {
		resp.Version = ver
	}
	return resp
}

func (r *Response) Errors(err ...Meta) *Response {
	r.Meta = err
	return r
}

func (r *Response) Success(code int) *Response {
	r.Meta = Meta{Code: StatusText(code)}
	return r
}

func (r *Response) Body(body interface{}) {
	r.Data = body
}

// APIStatusSuccess for standard request api status success
func (r *Response) APIStatusSuccess(w http.ResponseWriter, req *http.Request) *responseWriter {
	r.Success(StatusSuccess)
	return Status(w, req, StatusSuccess, r)
}

// APIStatusCreated for standard request api status created
func (r *Response) APIStatusCreated(w http.ResponseWriter, req *http.Request) *responseWriter {
	r.Success(StatusCreated)
	return Status(w, req, StatusCreated, r)
}

// APIStatusAccepted for standard request api status accepted
func (r *Response) APIStatusAccepted(w http.ResponseWriter, req *http.Request) *responseWriter {
	r.Success(StatusAccepted)
	return Status(w, req, StatusAccepted, r)
}

// APIStatusPermanentRedirect for standard request api status redirect
func (r *Response) APIStatusPermanentRedirect(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusPermanentRedirect),
		Type:    StatusCode(StatusPermanentRedirect),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusPermanentRedirect), err.Error()),
	})
	return Status(w, req, StatusPermanentRedirect, r)
}

// APIStatusBadRequest for standard request api status bad request
func (r *Response) APIStatusBadRequest(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusBadRequest),
		Type:    StatusCode(StatusBadRequest),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusBadRequest), err.Error()),
	})
	return Status(w, req, StatusBadRequest, r)
}

// APIStatusUnauthorized for standard request api status unauthorized
func (r *Response) APIStatusUnauthorized(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusUnauthorized),
		Type:    StatusCode(StatusUnauthorized),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusUnauthorized), err.Error()),
	})
	return Status(w, req, StatusUnauthorized, r)
}

// APIStatusPaymentRequired for standard request api status payment required
func (r *Response) APIStatusPaymentRequired(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusPaymentRequired),
		Type:    StatusCode(StatusPaymentRequired),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusPaymentRequired), err.Error()),
	})
	return Status(w, req, StatusPaymentRequired, r)
}

// APIStatusForbidden for standard request api status forbidden
func (r *Response) APIStatusForbidden(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusForbidden),
		Type:    StatusCode(StatusForbidden),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusForbidden), err.Error()),
	})
	return Status(w, req, StatusForbidden, r)
}

// APIStatusMethodNotAllowed for standard request api status not allowed
func (r *Response) APIStatusMethodNotAllowed(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusMethodNotAllowed),
		Type:    StatusCode(StatusMethodNotAllowed),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusMethodNotAllowed), err.Error()),
	})
	return Status(w, req, StatusMethodNotAllowed, r)
}

// APIStatusNotAcceptable for standard request api status StatusNotAcceptable
func (r *Response) APIStatusNotAcceptable(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusNotAcceptable),
		Type:    StatusCode(StatusNotAcceptable),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusNotAcceptable), err.Error()),
	})
	return Status(w, req, StatusNotAcceptable, r)
}

// APIStatusInvalidAuthentication for standard request api status StatusInvalidAuthentication
func (r *Response) APIStatusInvalidAuthentication(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusInvalidAuthentication),
		Type:    StatusCode(StatusInvalidAuthentication),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusInvalidAuthentication), err.Error()),
	})
	return Status(w, req, StatusInvalidAuthentication, r)
}

// APIStatusRequestTimeout for standard request api status StatusRequestTimeout
func (r *Response) APIStatusRequestTimeout(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusRequestTimeout),
		Type:    StatusCode(StatusRequestTimeout),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusRequestTimeout), err.Error()),
	})
	return Status(w, req, StatusRequestTimeout, r)
}

// APIStatusUnsupportedMediaType for standard request api status StatusUnsupportedMediaType
func (r *Response) APIStatusUnsupportedMediaType(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    StatusCode(StatusUnsupportedMediaType),
		Type:    StatusCode(StatusUnsupportedMediaType),
		Message: err.Error(),
	})
	return Status(w, req, StatusUnsupportedMediaType, r)
}

// APIStatusUnProcess for standard request api status StatusUnProcess
func (r *Response) APIStatusUnProcess(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusUnProcess),
		Type:    StatusCode(StatusUnProcess),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusUnProcess), err.Error()),
	})
	return Status(w, req, StatusUnProcess, r)
}

// APIStatusInternalError for standard request api status StatusInternalError
func (r *Response) APIStatusInternalError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusInternalError),
		Type:    StatusCode(StatusInternalError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusInternalError), err.Error()),
	})
	return Status(w, req, StatusInternalError, r)
}

// APIStatusBadGatewayError for standard request api status StatusBadGatewayError
func (r *Response) APIStatusBadGatewayError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusBadGatewayError),
		Type:    StatusCode(StatusBadGatewayError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusBadGatewayError), err.Error()),
	})
	return Status(w, req, StatusBadGatewayError, r)
}

// APIStatusServiceUnavailableError for standard request api status success
func (r *Response) APIStatusServiceUnavailableError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusServiceUnavailableError),
		Type:    StatusCode(StatusServiceUnavailableError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusServiceUnavailableError), err.Error()),
	})
	return Status(w, req, StatusServiceUnavailableError, r)
}

// APIStatusGatewayTimeoutError for standard request api status StatusGatewayTimeoutError
func (r *Response) APIStatusGatewayTimeoutError(w http.ResponseWriter, req *http.Request, err error) *responseWriter {
	r.Errors(Meta{
		Code:    strconv.Itoa(StatusGatewayTimeoutError),
		Type:    StatusCode(StatusGatewayTimeoutError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusGatewayTimeoutError), err.Error()),
	})
	return Status(w, req, StatusGatewayTimeoutError, r)
}

type responseWriter struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Response *Response
}

func Status(w http.ResponseWriter, r *http.Request, status int, v *Response) *responseWriter {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxResponse, status))
	return &responseWriter{
		Request:  r,
		Writer:   w,
		Response: v,
	}
}

func SemanticVersion(r *http.Request, label string, version string) {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxVersion, Version{
		Label:  label,
		Number: version,
	}))
}
