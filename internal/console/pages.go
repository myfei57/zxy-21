package console

import (
	"embed"
	"net/http"
)

//go:embed pages/*.html
var pagesFS embed.FS

// Page returns the bytes of one embedded console page.
func Page(name string) ([]byte, error) {
	return pagesFS.ReadFile("pages/" + name)
}

func pageHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := Page(name)
		if err != nil {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	}
}
