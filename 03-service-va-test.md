> [← 02 — Model và repository](02-model-va-repository.md) · [04 — Handler và router Gin →](04-handler-va-router.md)

# Viết test service trước rồi cài quy tắc nghiệp vụ

**Tạo các tệp:** internal/todo/service_test.go, internal/todo/service.go.

- [ ] **Viết test đỏ trước**

Tạo internal/todo/service_test.go:

~~~go
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
~~~

- [ ] **Chạy test để thấy lỗi**

~~~bash
go test ./internal/todo
~~~

Kết quả mong đợi: biên dịch thất bại vì NewService và ErrInvalidInput chưa tồn tại.

- [ ] **Cài service tối thiểu để test xanh**

Tạo internal/todo/service.go:

~~~go
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
~~~

- [ ] **Định dạng, chạy test, và commit**

~~~bash
gofmt -w internal/todo/service.go internal/todo/service_test.go
go test ./...
git add internal/todo
git commit -m "feat: add todo service"
~~~

Service biết các quy tắc như title không rỗng. Nó không biết JSON, HTTP hay PostgreSQL; nhờ Repository interface, test không cần Docker.

---

Tiếp theo: [04 — Handler và router Gin →](04-handler-va-router.md)
