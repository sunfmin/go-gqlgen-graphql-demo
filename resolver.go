package graphqldemo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		ID:     uuid.New().String(),
		Text:   "learn go",
		Done:   true,
		UserID: "1",
	},
	{
		ID:     uuid.New().String(),
		Text:   "learn graphql",
		Done:   false,
		UserID: "2",
	},
	{
		ID:     uuid.New().String(),
		Text:   "learn graphql in Go",
		Done:   false,
		UserID: "2",
	},
}

var storeUsers = []*api.User{
	{
		ID:   "1",
		Name: "john",
	},
	{
		ID:   "2",
		Name: "felix",
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

func (r *queryResolver) Todos(ctx context.Context, text *string) (todos []*api.Todo, err error) {
	fmt.Println("Calling Todos resolver")
	if text == nil {
		return storeTodos, nil
	}
	for _, t := range storeTodos {
		if strings.Index(strings.ToLower(t.Text), strings.ToLower(*text)) >= 0 {
			todos = append(todos, t)
		}
	}
	return
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *api.Todo) (*api.User, error) {
	return ctx.Value(userLoaderKey).(*UserLoader).Load(obj.UserID)
}

const userLoaderKey = "userloader"

func UserLoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := NewUserLoader(UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []string) (users []*api.User, errs []error) {
				fmt.Println("UserLoader Fetch", ids)
				for _, id := range ids {
					for _, u := range storeUsers {
						if u.ID == id {
							users = append(users, u)
						}
					}
				}
				return
			},
		})
		ctx := context.WithValue(r.Context(), userLoaderKey, userloader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
