package graphqldemo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/sunfmin/go-gqlgen-graphql-demo/api"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

var storeTodos = []*api.Todo{
	{
		ID:   uuid.New().String(),
		Text: "learn go",
		Done: true,
	},
	{
		ID:   uuid.New().String(),
		Text: "learn graphql",
		Done: false,
	},
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateTodo(ctx context.Context, input api.NewTodo) (*api.Todo, error) {
	fmt.Println("Calling CreateTodo", input)

	todo := &api.Todo{
		Text: input.Text,
		ID:   uuid.New().String(),
	}
	storeTodos = append(storeTodos, todo)
	return todo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*api.Todo, error) {
	fmt.Println("Calling Todos resolver")
	return storeTodos, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *api.Todo) (*api.User, error) {
	fmt.Println("Calling User resolver", obj)
	return &api.User{
		ID:   "123",
		Name: "Felix",
	}, nil
}
