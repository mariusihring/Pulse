package services

import (
	"context"
	"fmt"
	"pulse/graph/graphql_model"
	"pulse/internal/auth"
	"pulse/internal/db/models"
	"pulse/internal/validator"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	jwt *auth.JWT
}

func NewAuthService(db *gorm.DB, jwt *auth.JWT) *AuthService {
	return &AuthService{db: db, jwt: jwt}
}

func (s *AuthService) Register(ctx context.Context, input graphql_model.RegisterInput) (*graphql_model.AuthResponse, error) {
	if err := validator.ValidateRegisterInput(input); err != nil {
		return nil, err
	}

	var count int64

	if err := s.db.Model(&models.User{}).Where("email = ?", input.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("user with email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Assign default role
		defaultRole := &models.Role{}
		if err := tx.Where("name = ?", "USER").First(defaultRole).Error; err != nil {
			return err
		}

		if err := tx.Model(user).Association("Roles").Append(defaultRole); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.jwt.Generate(user.ID, []string{"USER"})
	if err != nil {
		return nil, err
	}

	return &graphql_model.AuthResponse{
		Token: token,
		User:  &graphql_model.User{ID: user.ID, Name: user.Name, Email: user.Email},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input graphql_model.LoginInput) (*graphql_model.AuthResponse, error) {
	user := &models.User{}
	if err := s.db.Preload("Roles").Where("email = ?", input.Email).First(user).Error; err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	token, err := s.jwt.Generate(user.ID, roles)
	if err != nil {
		return nil, err
	}

	return &graphql_model.AuthResponse{
		Token: token,
		User:  &graphql_model.User{ID: user.ID, Name: user.Name, Email: user.Email},
	}, nil
}
