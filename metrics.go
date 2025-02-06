package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Store(cfg.fileServerHits.Add(1))
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits.Load())
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	template := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(template))
	// tmpl, err := template.ParseFiles("admin.html")
	// if err != nil {
	// 	http.Error(w, "Error loading template", http.StatusInternalServerError)
	// 	return
	// }

	// tmpl.Execute(w, cfg.fileServerHits.Load())
}
