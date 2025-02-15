package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"fmt"
	graphql_model1 "pulse/graph/graphql_model"

	"github.com/google/uuid"
)

// Chains is the resolver for the chains field.
func (r *queryResolver) Chains(ctx context.Context) ([]*graphql_model1.Chain, error) {
	panic(fmt.Errorf("not implemented: Chains - chains"))
}

// Chain is the resolver for the chain field.
func (r *queryResolver) Chain(ctx context.Context, id uuid.UUID) (*graphql_model1.Chain, error) {
	panic(fmt.Errorf("not implemented: Chain - chain"))
}
