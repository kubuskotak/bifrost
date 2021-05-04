package bifrost

import (
	"context"
	"net/http"
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

func SemanticVersion(r *http.Request, label string, version string) {
	*r = *r.WithContext(context.WithValue(r.Context(), CtxVersion, Version{
		Label:  label,
		Number: version,
	}))
}
