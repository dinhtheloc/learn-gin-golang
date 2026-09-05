package todo

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidInput = errors.New("invalid todo input")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, input CreateTodoInput) (*Todo, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, ErrInvalidInput
	}

	item := &Todo{
		Title:       title,
		Description: input.Description,
		Completed:   false,
	}
	if err := s.repository.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) List(ctx context.Context) ([]Todo, error) {
	return s.repository.List(ctx)
}

func (s *Service) Get(ctx context.Context, id uint) (*Todo, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uint, input UpdateTodoInput) (*Todo, error) {
	if input.Title == nil && input.Description == nil && input.Completed == nil {
		return nil, ErrInvalidInput
	}

	item, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, ErrInvalidInput
		}
		item.Title = title
	}
	if input.Description != nil {
		item.Description = *input.Description
	}
	if input.Completed != nil {
		item.Completed = *input.Completed
	}

	if err := s.repository.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}
