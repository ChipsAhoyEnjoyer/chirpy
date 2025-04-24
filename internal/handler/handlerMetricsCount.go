package handler

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandlerMetricsCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	body := fmt.Sprintf(
		"<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>",
		cfg.FileserverHits.Load(),
	)
	w.Write([]byte(body))
}
