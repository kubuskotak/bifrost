package bifrost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemanticVersion(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(StatusSuccess) // set header code
	if got, want := w.Code, StatusSuccess; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
	SemanticVersion(r, "v1", "1.0.0")
	data := map[string]interface{}{
		"message": "transaksi telah sukses",
	}
	result := NewResponse(r)
	result.Body(data)
	result.APIStatusSuccess(w, r).WriteJSON()

	actual, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
	t.Log(string(actual))
	assert.Equal(t, result.Data, data)
}

func TestNew(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(StatusSuccess) // set header code
	if got, want := w.Code, StatusSuccess; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
	data := map[string]interface{}{
		"message": "transaksi telah sukses",
	}
	result := NewResponse(r)
	result.Body(data)
	result.APIStatusSuccess(w, r).WriteJSON()

	actual, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(actual))
	t.Log(result)
	assert.Equal(t, result.Data, data)
}

func TestResponseErrors(t *testing.T) {
	errs := make([]Meta, 0)
	errs = append(errs, Meta{
		Code:    strconv.Itoa(StatusInternalError),
		Type:    StatusCode(StatusInternalError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusInternalError), "constraint unique key duplicate"),
	})

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(StatusInternalError) // set header code
	if got, want := w.Code, StatusInternalError; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	result := NewResponse(r)
	result.APIStatusInternalError(w, r,
		fmt.Errorf("%s", "constraint unique key duplicate"),
	).WriteJSON()
	assert.Equal(t, result.Meta, errs)
	assert.Equal(t, strconv.Itoa(StatusInternalError), result.Meta.([]Meta)[0].Code)
	assert.Contains(t, result.Meta.([]Meta)[0].Message, "constraint unique key duplicate")
}

func TestResponseErrorsJSON(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	assert.NoError(t, err)
	w := httptest.NewRecorder()
	w.WriteHeader(StatusInternalError) // set header code
	if got, want := w.Code, StatusInternalError; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	errs := make([]Meta, 0)
	errs = append(errs, Meta{
		Code:    strconv.Itoa(StatusInternalError),
		Type:    StatusCode(StatusInternalError),
		Message: fmt.Sprintf("%s or %v", StatusText(StatusInternalError), "constraint unique key duplicate"),
	})
	result := NewResponse(r)
	result.Errors(errs...)
	result.APIStatusInternalError(w, r,
		fmt.Errorf("%s", "constraint unique key duplicate"),
	).WriteJSON()

	expected, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, result.Meta, errs)
	assert.Equal(t, strconv.Itoa(StatusInternalError), result.Meta.([]Meta)[0].Code)
	assert.Contains(t, result.Meta.([]Meta)[0].Message, "constraint unique key duplicate")
	assert.Equal(t, string(expected), strings.TrimSuffix(string(actual), "\n"))
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
	w.WriteHeader(StatusSuccess) // set header code
	if got, want := w.Code, StatusSuccess; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}

	result := NewResponse(r)
	result.APIStatusSuccess(w, r).WriteCSV(rows, "result-route-fleets") // Write http Body

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
