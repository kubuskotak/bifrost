package bifrost

import (
	"github.com/monoculum/formam"
	"net/http"
	"net/url"
	"strings"
)

func BindBody(r *http.Request, i interface{}) error {
	if r.ContentLength < 1 {
		return http.ErrContentLength
	}

	cType := r.Header.Get(HeaderContentType)
	switch {
	case strings.HasPrefix(cType, MIMEApplicationJSON):
		return RequestJSONBody(r, i)
	case strings.HasPrefix(cType, MIMEApplicationForm),
		strings.HasPrefix(cType, MIMEMultipartForm):
		p, err := params(r)
		if err != nil {
			return err
		}
		dec := formam.NewDecoder(
			&formam.DecoderOptions{TagName: "json"},
		)
		return dec.Decode(p, i)
	default:
		return http.ErrBodyNotAllowed
	}
}

func params(r *http.Request) (url.Values, error) {
	if strings.HasPrefix(r.Header.Get(HeaderContentType), MIMEMultipartForm) {
		if err := r.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
	}
	return r.Form, nil
}
