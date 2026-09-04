> [← 03 — Service và test](03-service-va-test.md) · [05 — Chạy API và kiểm thử →](05-chay-va-kiem-thu.md)

# Thêm Gin handler và router

**Tạo các tệp:** internal/todo/handler.go, internal/server/router.go.

### 4.1 handler.go

- [ ] Tạo internal/todo/handler.go:

~~~go
package todo

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

type errorResponse struct {
	Message string `json:"message"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return
	}

	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *Handler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid request body"})
		return
	}
	item, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "id must be a positive integer"})
		return 0, false
	}
	return uint(value), true
}

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		c.JSON(http.StatusBadRequest, errorResponse{Message: "invalid todo input"})
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Message: "todo not found"})
	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "internal server error"})
	}
}
~~~

ShouldBindJSON đọc body và áp dụng tag binding. Handler chỉ chuyển đổi HTTP thành lời gọi service, không chứa SQL.

### 4.2 router.go

- [ ] Tạo internal/server/router.go:

~~~go
package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"example.com/todo-api/internal/todo"
)

func NewRouter(handler *todo.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	todos := router.Group("/todos")
	{
		todos.POST("", handler.Create)
		todos.GET("", handler.List)
		todos.GET("/:id", handler.Get)
		todos.PATCH("/:id", handler.Update)
		todos.DELETE("/:id", handler.Delete)
	}
	return router
}
~~~

- [ ] **Định dạng, kiểm tra, commit**

~~~bash
gofmt -w internal/todo/handler.go internal/server/router.go
go test ./...
git add internal/todo/handler.go internal/server/router.go
git commit -m "feat: add Gin todo routes"
~~~

---

Tiếp theo: [05 — Chạy API và kiểm thử →](05-chay-va-kiem-thu.md)
