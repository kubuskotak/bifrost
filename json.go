package bifrost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONResponse set header content-type to json format
func JSONResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

// RequestJSONBody get json body format
func RequestJSONBody(r *http.Request, extract interface{}) error {
	// check method
	switch r.Method {
	case http.MethodPut:
	case http.MethodPost:
	case http.MethodGet:
		return fmt.Errorf("method is not allowed")
	}
	if r.Body == nil {
		return fmt.Errorf("there is no content")
	}
	if err := json.NewDecoder(r.Body).Decode(&extract); err != nil {
		return err
	}
	return nil
}

// ResponsePayload set payload for response http
func ResponsePayload(w http.ResponseWriter, r *http.Request, code int, payload interface{}) error {
	null := make(map[string]interface{})
	resp := &Response{
		Version: Version{
			Label:  "v1",
			Number: "0.1.0",
		},
		Meta:       Meta{Code: http.StatusText(code)},
		Data:       payload,
		Pagination: null,
	}
	if ver, ok := r.Context().Value(CtxVersion).(Version); ok {
		resp.Version = ver
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(resp); err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(code)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
