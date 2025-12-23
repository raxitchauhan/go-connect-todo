package handler

import (
	"context"
	"errors"

	todov1 "go-connect-todo/gen/todo/v1"
	"go-connect-todo/gen/todo/v1/todov1connect"
	"go-connect-todo/internal/repository"

	"connectrpc.com/connect"
)

type todo struct {
	todoRepo repository.TodoRepository
}

func NewTodo(repository repository.TodoRepository) todov1connect.TodoServiceHandler {
	return &todo{todoRepo: repository}
}

func (s *todo) List(ctx context.Context, req *todov1.ListRequest) (*todov1.ListResponse, error) {
	tasks, err := s.todoRepo.List(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &todov1.ListResponse{Tasks: tasks}, nil
}

func (s *todo) Complete(ctx context.Context, req *todov1.CompleteRequest) (*todov1.CompleteResponse, error) {
	task, err := s.todoRepo.Complete(ctx, req.Id)
	if err != nil {
		if errors.Is(err, repository.NotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &todov1.CompleteResponse{Task: task}, nil
}

func (s *todo) Add(ctx context.Context, req *todov1.AddRequest) (*todov1.AddResponse, error) {
	task, err := s.todoRepo.Add(ctx, req.Title)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &todov1.AddResponse{Task: task}, nil
}
