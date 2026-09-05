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
