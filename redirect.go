package bifrost

import (
	"net/http"
	"time"
)

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
