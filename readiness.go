package main

import "net/http"

// @Summary      Health check endpoint
// @Description  Returns 200 OK if the server is healthy
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /healthz [get]
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, map[string]string{
		"status": "ok",
	})
}

// func handlerReadiness(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(http.StatusText(http.StatusOK)))
// }
