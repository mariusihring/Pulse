package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"pulse/graph/generated"
	graphql_model1 "pulse/graph/graphql_model"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input graphql_model1.LoginInput) (*graphql_model1.AuthResponse, error) {
	return r.AuthService.Login(ctx, input)
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input graphql_model1.RegisterInput) (*graphql_model1.AuthResponse, error) {
	return r.AuthService.Register(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
