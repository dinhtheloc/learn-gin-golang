> [← 01 — Khởi tạo dự án và PostgreSQL](01-khoi-tao.md) · [03 — Service và test →](03-service-va-test.md)

# Định nghĩa TODO và lớp truy cập dữ liệu

**Tạo các tệp:** internal/todo/model.go, internal/todo/repository.go, internal/database/postgres.go.

### 2.1 model.go

- [ ] Tạo internal/todo/model.go:

~~~go
package todo

import "time"

type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTodoInput struct {
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description"`
}

// Con trỏ phân biệt field bị bỏ qua với completed: false.
type UpdateTodoInput struct {
	Title       *string `json:"title" binding:"omitempty,max=255"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}
~~~

Điểm cần nhớ: JSON tag quyết định tên thuộc tính API. Với PATCH, Completed phải là *bool: false là một giá trị hợp lệ, khác với không gửi field này.

### 2.2 repository.go

- [ ] Tạo internal/todo/repository.go:

~~~go
package todo

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("todo not found")

// Repository là hợp đồng. Service không biết và không cần biết GORM.
type Repository interface {
	Create(ctx context.Context, item *Todo) error
	List(ctx context.Context) ([]Todo, error)
	FindByID(ctx context.Context, id uint) (*Todo, error)
	Update(ctx context.Context, item *Todo) error
	Delete(ctx context.Context, id uint) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, item *Todo) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormRepository) List(ctx context.Context) ([]Todo, error) {
	var items []Todo
	err := r.db.WithContext(ctx).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *GormRepository) FindByID(ctx context.Context, id uint) (*Todo, error) {
	var item Todo
	err := r.db.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GormRepository) Update(ctx context.Context, item *Todo) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GormRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Todo{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
~~~

Repository là nơi duy nhất biết GORM. context.Context được Gin truyền vào để request bị hủy thì truy vấn cũng có thể dừng.

### 2.3 postgres.go

- [ ] Tạo internal/database/postgres.go:

~~~go
package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}
~~~

%w giữ lỗi gốc để sau này có thể kiểm tra bằng errors.Is. Timeout giúp API không treo vô hạn khi PostgreSQL không chạy.

- [ ] **Định dạng và kiểm tra biên dịch**

~~~bash
gofmt -w internal/todo/model.go internal/todo/repository.go internal/database/postgres.go
go test ./...
~~~

Kết quả mong đợi: mọi package biên dịch; chưa có test nào.

---

Tiếp theo: [03 — Service và test →](03-service-va-test.md)
