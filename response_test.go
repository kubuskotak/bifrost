package bifrost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupResponse() Response {
	return Response{
		Pagination: map[string]interface{}{},
	}
}

func TestSemanticVersion(t *testing.T) {
	response := setupResponse()
	label := "v1"
	version := "1.0.0"
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusOK) // set header code
	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	SemanticVersion(r, label, version)
	expected := map[string]interface{}{
		"message": "transaksi telah sukses",
	}

	JSONResponse(w)
	err = ResponsePayload(w, r, http.StatusOK, expected)
	assert.NoError(t, err)

	bytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	var actual Response
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err)

	response.Version = map[string]interface{}{
		"label":  label,
		"number": version,
	}
	response.Data = expected
	response.Meta = map[string]interface{}{"code": http.StatusText(http.StatusOK)}
	assert.Equal(t, response, actual)
}

func TestNew(t *testing.T) {
	response := setupResponse()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusOK) // set header code
	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
	expected := map[string]interface{}{
		"message": "transaksi telah sukses",
	}

	JSONResponse(w)
	err = ResponsePayload(w, r, http.StatusOK, expected)
	assert.NoError(t, err)

	bytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	var actual Response
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err)

	response.Version = map[string]interface{}{
		"label":  "v1",
		"number": "0.1.0",
	}
	response.Data = expected
	response.Meta = map[string]interface{}{"code": http.StatusText(http.StatusOK)}
	assert.Equal(t, response, actual)
}

func TestErrInternalServerError(t *testing.T) {
	//response := setupResponse()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	errReq := ErrInternalServerError(w, r, fmt.Errorf("%s or %v", http.StatusText(http.StatusInternalServerError), "constraint unique key duplicate"))
	assert.Error(t, errReq)
	if got, want := w.Code, http.StatusInternalServerError; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
}

func TestErrBadRequest(t *testing.T) {
	//response := setupResponse()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	errReq := ErrBadRequest(w, r, fmt.Errorf("%s or %v", http.StatusText(http.StatusBadRequest), "constraint unique key duplicate"))
	assert.Error(t, errReq)
	if got, want := w.Code, http.StatusBadRequest; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
}

func TestStatusBadGateway(t *testing.T) {
	//response := setupResponse()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	errReq := ErrBadGateway(w, r, fmt.Errorf("%s or %v", http.StatusText(http.StatusBadGateway), "constraint unique key duplicate"))
	assert.Error(t, errReq)
	if got, want := w.Code, http.StatusBadGateway; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
}

func TestStatusGatewayTimeout(t *testing.T) {
	//response := setupResponse()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	errReq := ErrGatewayTimeout(w, r, fmt.Errorf("%s or %v", http.StatusText(http.StatusGatewayTimeout), "constraint unique key duplicate"))
	assert.Error(t, errReq)
	if got, want := w.Code, http.StatusGatewayTimeout; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
}

func TestHttpErrStatus(t *testing.T) {
	status := []struct {
		label string
		code  int
		err   func(w http.ResponseWriter, r *http.Request, err error) error
	}{
		{http.StatusText(http.StatusBadRequest), http.StatusBadRequest, ErrBadRequest},
		{http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized, ErrUnauthorized},
		{http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired, ErrPaymentRequired},
		{http.StatusText(http.StatusForbidden), http.StatusForbidden, ErrForbidden},
		{http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed, ErrMethodNotAllowed},
		{http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable, ErrNotAcceptable},
		{http.StatusText(http.StatusProxyAuthRequired), http.StatusProxyAuthRequired, ErrProxyAuthRequired},
		{http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout, ErrRequestTimeout},
		{http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType, ErrUnsupportedMediaType},
		{http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity, ErrUnprocessableEntity},
		{http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError, ErrInternalServerError},
		{http.StatusText(http.StatusBadGateway), http.StatusBadGateway, ErrBadGateway},
		{http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable, ErrServiceUnavailable},
		{http.StatusText(http.StatusGatewayTimeout), http.StatusGatewayTimeout, ErrGatewayTimeout},
	}

	t.Run("http status", func(t *testing.T) {
		for _, tt := range status {
			t.Run(tt.label, func(t *testing.T) {
				r, err := http.NewRequest(http.MethodGet, "/", nil)
				assert.NoError(t, err)
				w := httptest.NewRecorder()
				errReq := tt.err(w, r, fmt.Errorf("%s or %v", tt.label, "constraint unique key duplicate"))
				assert.Error(t, errReq)
				if got, want := w.Code, tt.code; got != want {
					t.Fatalf("status code got: %d, want %d", got, want)
				}
			})
		}
	})
}

func TestResponseCSV(t *testing.T) {
	rows := make([][]string, 0)
	rows = append(rows, []string{"SO Number", "Nama Warung", "Area", "Fleet Number", "Jarak Warehouse", "Urutan"})
	rows = append(rows, []string{"SO45678", "WPD00011", "Jakarta Selatan", "1", "45.00", "1"})
	rows = append(rows, []string{"SO45645", "WPD001123", "Jakarta Selatan", "1", "43.00", "2"})
	rows = append(rows, []string{"SO45645", "WPD003343", "Jakarta Selatan", "1", "43.00", "3"})

	r, err := http.NewRequest(http.MethodGet, "/csv", nil)

	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(http.StatusOK) // set header code
	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	erCVS := ResponseCSVPayload(w, r, http.StatusOK, rows, "result-route-fleets")
	assert.NoError(t, erCVS)

	actual, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `SO Number,Nama Warung,Area,Fleet Number,Jarak Warehouse,Urutan
SO45678,WPD00011,Jakarta Selatan,1,45.00,1
SO45645,WPD001123,Jakarta Selatan,1,43.00,2
SO45645,WPD003343,Jakarta Selatan,1,43.00,3
`, string(actual))
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")

}
