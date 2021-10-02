package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var r R

type Todo struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Weight      int       `json:"weight"`
	Tags        []string  `json:"tags"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


type R interface {
	AddTodo(context.Context, Todo) (Todo, error)
	GetTodoByID(context.Context, string) (Todo, error)
	UpdateTodoStatus(context.Context, Todo, bool) (Todo, error)
	UpdateTodoTags(context.Context, Todo, []string) (Todo, error)
	ListTodo(ctx context.Context, sortBy string, desc bool, all bool) ([]Todo, error)
	DeleteTodo(context.Context, string) error
}

type contextKey string

const contextKeyR contextKey = "repository"

func SetR(repo R) {
	r = repo
}

func GetR() (repo R, e error) {
	if r == nil {
		return nil, errors.New("global repository is empty")
	}

	return r, nil
}
