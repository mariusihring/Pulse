package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"fmt"
	"pulse/graph/generated"
	"pulse/graph/graphql_model"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input graphql_model.LoginInput) (*graphql_model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: Login - login"))
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input graphql_model.RegisterInput) (*graphql_model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: Register - register"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
