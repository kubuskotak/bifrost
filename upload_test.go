package bifrost

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestUploadBase64Jpeg(t *testing.T) {
	f, err := os.Open("./ava.jpeg")
	assert.NoError(t, err)

	img64, err := ImgToBase64(f)
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	mime, err := Base64ToImg(buf, img64)
	assert.NoError(t, err)
	t.Log(mime)
}

func TestUploadBase64PNG(t *testing.T) {
	f, err := os.Open("./amazon-sns.png")
	assert.NoError(t, err)

	img64, err := ImgToBase64(f)
	assert.NoError(t, err)

	buf := &bytes.Buffer{}
	mime, err := Base64ToImg(buf, img64)
	assert.NoError(t, err)
	t.Log(mime)
}