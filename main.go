package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Brian202208/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB URL is not found in the environment")

	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeed)
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	v1Router.Post("/feeds_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feeds_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	v1Router.Delete("/feeds_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Sever starting on port %v", portString)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
