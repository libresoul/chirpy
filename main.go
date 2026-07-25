package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/libresoul/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	jwt_secret     string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if dbURL == "" {
		log.Fatal("DB_URL missing")
	}
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET missing")
	}

	db, _ := sql.Open("postgres", dbURL)
	err := db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	err = runMigrations("postgres", db, os.DirFS("sql/schema"))
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8000"

	apiCfg := &apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		jwt_secret:     jwtSecret,
	}

	mux := http.NewServeMux()

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerRediness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + "8080",
		Handler: mux,
	}

	log.Println("Server started on port", port)
	log.Fatal(server.ListenAndServe())
}
