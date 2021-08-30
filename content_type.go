package bifrost

import (
	"context"
	"net/http"
	"strings"
)

var (
	ContentTypeCtxKey = &ctxRender{HeaderContentType}
)

// ContentType is an enumeration of common HTTP content types.
type ContentType int

// ContentTypes handled by this package.
const (
	ContentTypeUnknown = iota
	ContentTypePlainText
	ContentTypeHTML
	ContentTypeJSON
	ContentTypeXML
	ContentTypeForm
	ContentTypeMultipartForm
	ContentTypeEventStream
)

// GetContentType is a middleware that forces response Content-Type.
func GetContentType(s string) ContentType {
	cType := strings.TrimSpace(strings.Split(s, ";")[0])
	switch {
	case strings.HasPrefix(cType, MIMETextHTML):
		return ContentTypeHTML
	case strings.HasPrefix(cType, MIMETextXML):
		return ContentTypeXML
	case strings.HasPrefix(cType, MIMEApplicationJSON):
		return ContentTypeJSON
	case strings.HasPrefix(cType, MIMEApplicationForm):
		return ContentTypeForm
	case strings.HasPrefix(cType, MIMEMultipartForm):
		return ContentTypeMultipartForm
	case strings.HasPrefix(cType, MIMETextPlain):
		return ContentTypePlainText
	case strings.HasPrefix(cType, MIMEOctetStream):
		return ContentTypeEventStream
	default:
		return ContentTypeUnknown
	}
}

// SetContentType is a middleware that forces response Content-Type.
func SetContentType(contentType ContentType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ContentTypeCtxKey, contentType))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// GetRequestContentType is a helper function that returns ContentType based on
// context or request headers.
func GetRequestContentType(r *http.Request) ContentType {
	if contentType, ok := r.Context().Value(ContentTypeCtxKey).(ContentType); ok {
		return contentType
	}
	return GetContentType(r.Header.Get("Content-Type"))
}

// GetIdxContentType is a middleware that forces response Content-Type.
func GetIdxContentType(s ContentType) string {
	switch s {
	case ContentTypeHTML:
		return MIMETextHTML
	case ContentTypeXML:
		return MIMETextXML
	case ContentTypeJSON:
		return MIMEApplicationJSON
	case ContentTypeForm:
		return MIMEApplicationForm
	case ContentTypeMultipartForm:
		return MIMEMultipartForm
	case ContentTypePlainText:
		return MIMETextPlain
	case ContentTypeEventStream:
		return MIMEOctetStream
	default:
		return ""
	}
}
