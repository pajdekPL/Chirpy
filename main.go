package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/PajdekPL/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"
	cfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.handlerGetChirp)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", cfg.handlerChangeEmailPassword)
	mux.HandleFunc("DELETE /api/chirps/{id}", cfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaWebhook)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
