package auth

import (
	"pulse/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"sub"`             // Supabase uses 'sub' for user ID
	Email  string    `json:"email"`           // Supabase email
	Role   string    `json:"role"`            // Supabase role
	Roles  []string  `json:"roles,omitempty"` // Your custom roles
}

type JWT struct {
	secretKey []byte
	duration  time.Duration
}

func New(cfg *config.Config) *JWT {
	return &JWT{
		secretKey: []byte(cfg.JWT.Secret), // This should be your Supabase JWT secret
		duration:  cfg.JWT.Duration,
	}
}

func (j *JWT) Validate(tokenStr string) (*Claims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	// Get user roles from database if needed
	// You might want to add a method to fetch roles here

	return &claims, nil
}
