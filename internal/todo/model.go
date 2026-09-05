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
