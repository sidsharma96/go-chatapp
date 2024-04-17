package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sidsharma96/chatapp-backend/db/sqlc"
)

type Server struct {
	DbQueries  *sqlc.Queries
	HttpServer *http.Server
	ChiRouter  *chi.Mux
}

func NewServer(port string) Server {
	godotenv.Load(".env")

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("Database url is not found in the environment!")
	}

	conn, dbErr := sql.Open("postgres", dbUrl)
	if dbErr != nil {
		log.Fatal("Cannot connect to database - ", dbErr)
	}

	driver, driverErr := postgres.WithInstance(conn, &postgres.Config{})
	if driverErr != nil {
		log.Fatal("Error opening postgres driver - ", driverErr)
	}
	migration, mErr := migrate.NewWithDatabaseInstance("file://db/schema", "postgres", driver)
	if mErr != nil {
		log.Fatal("Error opening migrate instance - ", mErr)
	}
	upErr := migration.Up()
	if upErr != nil && upErr != migrate.ErrNoChange {
		log.Fatal("Error migrating up - ", upErr)
	}
	queries := sqlc.New(conn)

	defaultRouter := chi.NewRouter()
	defaultRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "withCredentials"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	defaultRouter.Get("/", defaultHandler)
	defaultRouter.Get("/error", defaultErrorHandler)

	srv := &http.Server{
		Handler: defaultRouter,
		Addr:    ":" + port,
	}

	return Server{
		DbQueries:  queries,
		HttpServer: srv,
		ChiRouter:  defaultRouter,
	}
}

func (s *Server) Start() {
	Auth{}.router(s)
	User{}.router(s)

	pool := NewPool(s)
	go pool.Start()
	pool.router(s)

	fmt.Printf("Running server on port: %s", s.HttpServer.Addr)
	err := s.HttpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	type defaultResponse struct {
		Message string `json:"message"`
	}

	jsonHandler(w, http.StatusOK, defaultResponse{
		Message: "Hello there",
	})
}

func defaultErrorHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, http.StatusNotFound, "You shall not pass!")
}

func errorHandler(w http.ResponseWriter, status int, errorMsg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	jsonHandler(w, status, errorResponse{
		Error: errorMsg,
	})
}

func jsonHandler(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
