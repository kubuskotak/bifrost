package bifrost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// JSONResponse set header content-type to json format
func JSONResponse(w http.ResponseWriter) {
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
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

// ResponseJSONPayload set payload for response http
func ResponseJSONPayload(w http.ResponseWriter, r *http.Request, code int, responses ...interface{}) error {
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	null := make(map[string]interface{})
	resp := &Response{
		Version: Version{
			Label:  "v1",
			Number: "0.1.0",
		},
		Meta:       Meta{Code: http.StatusText(code)},
		Pagination: null,
	}
	if ver, ok := r.Context().Value(CtxVersion).(Version); ok {
		resp.Version = ver
	}
	data := make(map[string]interface{}, 0)
	for _, r := range responses {
		switch r.(type) {
		case Pagination:
			resp.Pagination = r
		case map[string]interface{}:
			for k, v := range r.(map[string]interface{}) {
				data[k] = v
			}
		default:
			iType := reflect.TypeOf(r)
			switch iType.Kind() {
			case reflect.Slice, reflect.Array:
				s := reflect.ValueOf(r)
				l := strings.Split(s.Index(0).Type().String(), ".")

				dataList := make([]map[string]interface{}, 0)
				for n := 0; n < s.Len(); n++ {
					b, err := json.Marshal(s.Index(n).Interface())
					if err != nil {
						w.Header().Set("X-Content-Type-Options", "nosniff")
						w.WriteHeader(http.StatusInternalServerError)
						return err
					}
					tempData := make(map[string]interface{}, 0)
					if err := json.Unmarshal(b, &tempData); err != nil {
						w.Header().Set("X-Content-Type-Options", "nosniff")
						w.WriteHeader(http.StatusInternalServerError)
						return err
					}
					dataList = append(dataList, tempData)
				}
				data[ToDelimited(l[len(l)-1], '_')] = dataList
			default:
				b, err := json.Marshal(r)
				if err != nil {
					w.Header().Set("X-Content-Type-Options", "nosniff")
					w.WriteHeader(http.StatusInternalServerError)
					return err
				}
				tempData := make(map[string]interface{}, 0)
				if err := json.Unmarshal(b, &tempData); err != nil {
					w.Header().Set("X-Content-Type-Options", "nosniff")
					w.WriteHeader(http.StatusInternalServerError)
					return err
				}
				for k, v := range tempData {
					data[k] = v
				}
			}

		}
	}
	resp.Data = data

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
