package main

import (
	"net/http"
	"os"
	"pulse/graph/generated"
	"pulse/graph/resolvers"
	"pulse/internal/auth"
	"pulse/internal/config"
	"pulse/internal/db"
	"pulse/internal/services"

	"github.com/charmbracelet/log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const defaultPort = "8080"

func main() {
	cfg := config.Load()

	jwtAuth := auth.New(cfg)
	database := db.InitDB(cfg)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Content-Type"))

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			AuthService: services.NewAuthService(database, jwtAuth),
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("Pulse", "/query"))
	router.Handle("/query", srv)

	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Infof("Connect to http://localhost:%s/ for Graphql Playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
