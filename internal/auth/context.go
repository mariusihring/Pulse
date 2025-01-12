package auth

import (
	"context"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id" // Match the key from middleware
	RolesKey  contextKey = "roles"   // Match the key from middleware
)

func UserFromContext(ctx context.Context) *uuid.UUID {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return nil
	}
	if userID, ok := val.(uuid.UUID); ok {
		return &userID
	}
	return nil
}

func RolesFromContext(ctx context.Context) []string {
	val := ctx.Value(RolesKey)
	if val == nil {
		return nil
	}
	if roles, ok := val.([]string); ok {
		return roles
	}
	return nil
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithRoles adds roles to context
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, RolesKey, roles)
}
