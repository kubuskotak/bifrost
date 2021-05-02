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

const dom string = "__CONSTANTA__"

var (
	RenderContext = ctxRender{Name: "context render"}
)

func RenderWriter(r *http.Request, value interface{}) {
	constant := map[string]interface{}{
		dom: value,
	}
	*r = *r.WithContext(context.WithValue(r.Context(), RenderContext, constant))
}

func RenderReader(r *http.Request) interface{} {
	return r.Context().Value(RenderContext)
}
