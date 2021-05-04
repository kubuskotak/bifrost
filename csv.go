package bifrost

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
)

// ResponseCSVPayload set csv payload for response http
func ResponseCSVPayload(w http.ResponseWriter, r *http.Request, code int, rows [][]string, filename string) error {
	buf := &bytes.Buffer{}
	xCsv := csv.NewWriter(buf)

	for _, row := range rows {
		if err := xCsv.Write(row); err != nil {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
	}
	xCsv.Flush()

	if err := xCsv.Error(); err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Description", "File Transfer")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")

	w.WriteHeader(code)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
