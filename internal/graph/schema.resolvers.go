package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/db"
)

func (r *installResolver) ID(ctx context.Context, obj *db.Install) (uuid.UUID, error) {
	return obj.InstallID, nil
}

func (r *installResolver) Username(ctx context.Context, obj *db.Install) (string, error) {
	return obj.SshUsername, nil
}

func (r *installResolver) Hostname(ctx context.Context, obj *db.Install) (string, error) {
	return obj.SshHostname, nil
}

func (r *installResolver) Port(ctx context.Context, obj *db.Install) (int, error) {
	return int(obj.SshPort), nil
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*db.Install, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Installs(ctx context.Context) ([]db.Install, error) {
	installs, err := r.Data.GetAllInstalls(ctx)
	return installs, err
}

// Install returns InstallResolver implementation.
func (r *Resolver) Install() InstallResolver { return &installResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type installResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
