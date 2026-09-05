package todo

import (
	"context"
	"errors"
	"testing"
)

type memoryRepository struct {
	nextID uint
	items  map[uint]*Todo
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{items: map[uint]*Todo{}}
}

func copyTodo(item *Todo) *Todo {
	value := *item
	return &value
}

func (r *memoryRepository) Create(_ context.Context, item *Todo) error {
	r.nextID++
	item.ID = r.nextID
	r.items[item.ID] = copyTodo(item)
	return nil
}

func (r *memoryRepository) List(_ context.Context) ([]Todo, error) {
	result := make([]Todo, 0, len(r.items))
	for id := uint(1); id <= r.nextID; id++ {
		if item, ok := r.items[id]; ok {
			result = append(result, *copyTodo(item))
		}
	}
	return result, nil
}

func (r *memoryRepository) FindByID(_ context.Context, id uint) (*Todo, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return copyTodo(item), nil
}

func (r *memoryRepository) Update(_ context.Context, item *Todo) error {
	if _, ok := r.items[item.ID]; !ok {
		return ErrNotFound
	}
	r.items[item.ID] = copyTodo(item)
	return nil
}

func (r *memoryRepository) Delete(_ context.Context, id uint) error {
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func TestCreateTrimsTitleAndStartsIncomplete(t *testing.T) {
	service := NewService(newMemoryRepository())

	item, err := service.Create(context.Background(), CreateTodoInput{
		Title: "  Learn Gin  ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if item.ID != 1 || item.Title != "Learn Gin" || item.Completed {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestCreateRejectsBlankTitle(t *testing.T) {
	service := NewService(newMemoryRepository())

	_, err := service.Create(context.Background(), CreateTodoInput{Title: "   "})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Create() error = %v, want ErrInvalidInput", err)
	}
}

func TestUpdateAcceptsExplicitCompletedValue(t *testing.T) {
	service := NewService(newMemoryRepository())
	item, err := service.Create(context.Background(), CreateTodoInput{Title: "Read Go tour"})
	if err != nil {
		t.Fatal(err)
	}

	completed := true
	updated, err := service.Update(context.Background(), item.ID, UpdateTodoInput{Completed: &completed})
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Completed {
		t.Fatal("Completed = false, want true")
	}
}
