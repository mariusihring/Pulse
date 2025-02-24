package main

import (
	"net/http"
	"pulse/graph/generated"
	"pulse/graph/resolvers"
	pulse_middleware "pulse/internal/auth/middleware"
	"pulse/internal/db"
	"pulse/internal/pubsub"

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

	database := db.InitDB()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization"))

	// Create PubSUb
	redisAddr := "localhost:6379"
	redisPassword := ""
	redisDB := 0

	// Create Services

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			RedisPubSub: pubsub.NewRedisPubSub(redisAddr, redisPassword, redisDB),
		},
		Directives: generated.DirectiveRoot{
			Auth: pulse_middleware.AuthDirective,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{})

	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("Pulse", "/query"))
	router.Handle("/query", pulse_middleware.Auth(database)(srv))

	/*
		router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	*/

	// log.Infof("Connect to http://localhost:%s/ for Graphql Playground", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":3001", router))
}
