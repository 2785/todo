package repositories

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/google/uuid"
	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"golang.org/x/sync/errgroup"
)

func NewFlatFile(dir string) (*FlatFile, error) {
	db, err := scribble.New(dir, &scribble.Options{})
	if err != nil {
		return nil, err
	}

	return &FlatFile{db}, nil
}

type FlatFile struct {
	d *scribble.Driver
}

func (ff *FlatFile) AddTodo(ctx context.Context, todo Todo) (created Todo, e error) {
	if todo.ID == uuid.Nil {
		todo.ID = uuid.New()
	}
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	ech := make(chan error)

	go func() {
		ech <- ff.d.Write("todos", todo.ID.String(), todo)
	}()

	select {
	case err := <-ech:
		if err != nil {
			e = err
			return
		}
	case <-ctx.Done():
		e = errors.New("context cancelled")
		return
	}

	return todo, nil
}

func (ff *FlatFile) GetTodoByID(ctx context.Context, id string) (todo Todo, e error) {
	ech := make(chan error)

	go func() {
		ech <- ff.d.Read("todos", id, &todo)
	}()

	select {
	case err := <-ech:
		if err != nil {
			e = err
		}
	case <-ctx.Done():
		e = errors.New("context cancelled")
	}

	return
}

func (ff *FlatFile) UpdateTodoStatus(ctx context.Context, todo Todo, newStatus bool) (updated Todo, e error) {
	ech := make(chan error)

	go func() {
		if err := ff.d.Read("todos", todo.ID.String(), &updated); err != nil {
			ech <- err
		}
		updated.UpdatedAt = time.Now()
		updated.Done = newStatus
		ech <- ff.d.Write("todos", updated.ID.String(), &updated)
	}()

	select {
	case err := <-ech:
		if err != nil {
			e = err
		}
	case <-ctx.Done():
		e = errors.New("context cancelled")
	}

	return
}

func (ff *FlatFile) UpdateTodoTags(ctx context.Context, todo Todo, tags []string) (updated Todo, e error) {
	ech := make(chan error)

	go func() {
		if err := ff.d.Read("todos", todo.ID.String(), &updated); err != nil {
			ech <- err
		}
		updated.UpdatedAt = time.Now()
		updated.Tags = tags
		ech <- ff.d.Write("todos", updated.ID.String(), &updated)
	}()

	select {
	case err := <-ech:
		if err != nil {
			e = err
		}
	case <-ctx.Done():
		e = errors.New("context cancelled")
	}

	return
}

func (ff *FlatFile) ListTodo(ctx context.Context, by string, desc bool, all bool) (res []Todo, e error) {
	switch by {
	case "weight":
	case "created_at":
	case "updated_at":
	default:
		e = errors.Errorf("unknown field %q to sort by", by)
	}

	allTodos, err := ff.GetAll(ctx)
	if err != nil {
		e = err
		return
	}

	allTodosPtr := make([]*Todo, len(allTodos))
	for i := range allTodos {
		allTodosPtr[i] = &allTodos[i]
	}

	if !all {
		allTodosPtr = funk.Filter(allTodosPtr, func(t *Todo) bool { return !t.Done }).([]*Todo)
	}

	switch by {
	case "weight":
		sort.SliceStable(allTodosPtr, func(i, j int) bool {
			if desc {
				return allTodosPtr[i].Weight > allTodosPtr[j].Weight
			}
			return allTodosPtr[i].Weight < allTodosPtr[j].Weight
		})
	case "created_at":
		sort.SliceStable(allTodosPtr, func(i, j int) bool {
			if desc {
				return allTodosPtr[i].CreatedAt.After(allTodosPtr[j].CreatedAt)
			}
			return allTodosPtr[i].CreatedAt.Before(allTodosPtr[j].CreatedAt)
		})
	case "updated_at":
		sort.SliceStable(allTodosPtr, func(i, j int) bool {
			if desc {
				return allTodosPtr[i].CreatedAt.After(allTodosPtr[j].CreatedAt)
			}
			return allTodosPtr[i].CreatedAt.Before(allTodosPtr[j].CreatedAt)
		})
	}

	res = make([]Todo, len(allTodosPtr))
	for i := range allTodosPtr {
		res[i] = *allTodosPtr[i]
	}

	return
}

func (ff *FlatFile) DeleteTodo(ctx context.Context, id string) error {
	ech := make(chan error)

	go func() {
		ech <- ff.d.Delete("todos", id)
	}()

	select {
	case err := <-ech:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return errors.New("context cancelled")
	}

	return nil
}

func (ff *FlatFile) GetAll(ctx context.Context) (res []Todo, e error) {
	ech := make(chan error)

	go func() {
		strings, err := ff.d.ReadAll("todos")
		if err != nil {
			ech <- err
		}
		res = make([]Todo, len(strings))

		eg := &errgroup.Group{}
		for i, v := range strings {
			ind, content := i, v
			eg.Go(func() error {
				todo := &Todo{}
				err := json.Unmarshal([]byte(content), todo)
				if err != nil {
					return err
				}
				res[ind] = *todo
				return nil
			})
		}

		ech <- eg.Wait()
	}()

	select {
	case err := <-ech:
		if err != nil {
			e = err
		}
	case <-ctx.Done():
		e = errors.New("context cancelled")
	}

	return
}
