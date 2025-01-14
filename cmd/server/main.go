package main

import (
	"net/http"
	"pulse/graph/generated"
	"pulse/graph/resolvers"
	"pulse/internal/auth"
	pulse_middleware "pulse/internal/auth/middleware"
	"pulse/internal/config"
	"pulse/internal/db"
	"pulse/internal/services"
	"pulse/internal/services/loaders"

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

	// Create Loaders
	solana_loader := loaders.NewSolanaLoader(cfg)

	// Create Services

	wallet_service := services.NewWalletService(database, solana_loader)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{
			AuthService:   services.NewAuthService(database, jwtAuth),
			WalletService: wallet_service,
		},
		Directives: generated.DirectiveRoot{
			Auth: pulse_middleware.AuthDirective,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.POST{})

	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("Pulse", "/query"))
	router.Handle("/query", pulse_middleware.Auth(jwtAuth)(srv))

	/*
		router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	*/

	log.Infof("Connect to http://localhost:%s/ for Graphql Playground", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(cfg.Server.Host+cfg.Server.Port, router))
}
