package bifrost

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// WriteJSON Write writes the data to http Response writer
func (r *responseWriter) WriteJSON() {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(r.Response); err != nil {
		http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Request.Context().Value(CtxResponse).(int); ok {
		r.Writer.WriteHeader(status)
	}
	_, err := r.Writer.Write(buf.Bytes())
	if err != nil {
		http.Error(r.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
