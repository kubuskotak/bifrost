package bifrost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"testing"
)

type person struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func TestBindBodyJSON(t *testing.T) {
	body := person{
		Name:    "surya",
		Address: "kencana mukti 21",
	}

	b, err := json.Marshal(body)
	assert.NoError(t, err)
	r, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	assert.NoError(t, err)
	r.Header.Set(HeaderContentType, MIMEApplicationJSON)

	var p person
	err = BindBody(r, &p)
	assert.NoError(t, err)
	assert.Equal(t, body, p)
}

func TestBindBodyForm(t *testing.T) {
	p := person{
		Name:    "surya",
		Address: "kencana mukti 21",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	contentType := fmt.Sprintf("%s; boundary=%q",
		MIMEMultipartForm, writer.Boundary())
	err := writer.WriteField("name", p.Name)
	assert.NoError(t, err)
	err = writer.WriteField("address", p.Address)
	assert.NoError(t, err)

	_ = writer.Close()
	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.NoError(t, err)
	r.Header.Set(HeaderContentType, contentType)

	var pn person
	err = BindBody(r, &pn)
	assert.NoError(t, err)
	assert.Equal(t, p, pn)
}

func TestBindBodyFailContentType(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	contentType := fmt.Sprintf("%s; boundary=%q",
		MIMETextHTML, writer.Boundary())

	_ = writer.Close()
	r, err := http.NewRequest(http.MethodPost, "/", body)
	assert.NoError(t, err)
	r.Header.Set(HeaderContentType, contentType)

	var pn person
	err = BindBody(r, &pn)
	assert.Error(t, err)
}
