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
