package bifrost

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func ImgToBase64(r io.ReadCloser) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	var base64Img string
	content := base64.StdEncoding.EncodeToString(b)
	mimeType := http.DetectContentType(b)
	switch mimeType {
	case MIMEImageJPEG:
		base64Img = fmt.Sprintf("data:%s;base64,%v", MIMEImageJPEG, content)
	case MIMEImagePNG:
		base64Img = fmt.Sprintf("data:%s;base64,%v", MIMEImagePNG, content)
	}
	return base64Img, nil
}

func Base64ToImg(w io.Writer, base string) (string, error) {
	var (
		content  string
		mimeType string
	)

	splitter := strings.Split(base, ",")
	if len(splitter) > 1 {
		content = splitter[1]
	} else {
		content = splitter[0]
	}
	switch {
	case strings.Contains(base, MIMEImageJPEG):
		mimeType = MIMEImageJPEG
	case strings.Contains(base, MIMEImagePNG):
		mimeType = MIMEImagePNG
	default:
		mimeType = "unknown"
	}
	b, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return mimeType, err
	}
	if _, err = w.Write(b); err != nil {
		return mimeType, err
	}
	return mimeType, nil
}
