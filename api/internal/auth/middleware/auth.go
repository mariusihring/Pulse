package middleware

import (
	"context"
	"fmt"
	"net/http"
	"pulse/internal/auth"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

func Auth(jwt *auth.JWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwt.Validate(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Use the same keys from auth package
			ctx := context.WithValue(r.Context(), auth.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, auth.RolesKey, claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	userId := ctx.Value(auth.UserIDKey)
	if userId == nil {
		return nil, fmt.Errorf("access denied")
	}
	return next(ctx)
}

func HasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, role string) (interface{}, error) {
	roles, ok := ctx.Value(auth.RolesKey).([]string)
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