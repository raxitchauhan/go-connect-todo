package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	todov1 "go-connect-todo/gen/todo/v1"
)

var NotFound = errors.New("not_found")

type TodoRepository interface {
	List(ctx context.Context) ([]*todov1.Task, error)
	Complete(ctx context.Context, id string) (*todov1.Task, error)
	Add(ctx context.Context, title string) (*todov1.Task, error)
}

type todo struct {
	counter atomic.Int64
	mu      *sync.Mutex
	tasks   []*todov1.Task
}

func New() TodoRepository {
	t := &todo{
		mu: &sync.Mutex{},
		tasks: []*todov1.Task{
			{Id: "1", Title: "Buy groceries", Completed: false},
			{Id: "2", Title: "Walk the dog", Completed: false},
		},
	}

	t.counter.Add(2)

	return t
}

func (s *todo) List(_ context.Context) ([]*todov1.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.tasks, nil
}

func (s *todo) Complete(_ context.Context, id string) (*todov1.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].Id == id {
			s.tasks[i].Completed = true

			return s.tasks[i], nil
		}
	}

	return nil, NotFound
}

func (s *todo) Add(_ context.Context, title string) (*todov1.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.counter.Add(1)
	task := &todov1.Task{
		Id:        fmt.Sprintf("%d", id),
		Title:     title,
		Completed: false,
	}

	s.tasks = append(s.tasks, task)

	return task, nil
}
