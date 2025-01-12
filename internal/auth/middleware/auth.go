package middleware

import (
	"context"
	"fmt"
	"net/http"
	"pulse/internal/auth"

	"github.com/99designs/gqlgen/graphql"
)

type contextKey string

const (
	UserIdKey contextKey = "user_id"
	RolesKey  contextKey = "roles"
)

func Auth(jwt *auth.JWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwt.Validate(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserIdKey, claims.UserID)
			ctx = context.WithValue(ctx, RolesKey, claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	userId := ctx.Value(UserIdKey)
	if userId == nil {
		return nil, fmt.Errorf("access denied")
	}
	return next(ctx)
}

func HasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role string) (interface{}, error) {
	roles, ok := ctx.Value(RolesKey).([]string)
	if !ok {
		return nil, fmt.Errorf("access denied")
	}

	for _, r := range roles {
		if r == role {
			return next(ctx)
		}
	}

	return nil, fmt.Errorf("insufficient permissions")
}
