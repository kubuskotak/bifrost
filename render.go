package bifrost

import (
	"context"
	"net/http"
)

type ctxRender struct {
	Name string
}

func (r *ctxRender) String() string {
	return "context value " + r.Name
}

var (
	RenderContext = ctxRender{Name: "context render"}
)

func RenderWriter(r *http.Request, value interface{}) {
	*r = *r.WithContext(context.WithValue(r.Context(), RenderContext, value))
}

func RenderReader(r *http.Request) interface{} {
	return r.Context().Value(RenderContext)
}
