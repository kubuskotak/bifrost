package bifrost

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	assert.NoError(t, err)
	w := httptest.NewRecorder()
	Redirect(w, r, "/redirect")
	if got, want := w.Code, http.StatusMovedPermanently; got != want {
		t.Fatalf("status code got: %d, want %d", got, want)
	}
}
