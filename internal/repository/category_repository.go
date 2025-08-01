package repository

import (
	"FinalTask/config"
	"FinalTask/internal/models"
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, cat *models.Category) error
	List(ctx context.Context) ([]*models.Category, error)
	FindByID(ctx context.Context, id uint) (*models.Category, error)
	Update(ctx context.Context, cat *models.Category) error
	Delete(ctx context.Context, id uint) error
}

type categoryRepo struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepo{}
}

func (r *categoryRepo) Create(ctx context.Context, cat *models.Category) error {
	return config.DB.WithContext(ctx).Create(cat).Error
}

func (r *categoryRepo) List(ctx context.Context) ([]*models.Category, error) {
	var list []*models.Category
	err := config.DB.WithContext(ctx).Find(&list).Error
	return list, err
}

func (r *categoryRepo) FindByID(ctx context.Context, id uint) (*models.Category, error) {
	var cat models.Category
	err := config.DB.WithContext(ctx).First(&cat, id).Error
	return &cat, err
}

func (r *categoryRepo) Update(ctx context.Context, cat *models.Category) error {
	return config.DB.WithContext(ctx).Save(cat).Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Delete(&models.Category{}, id).Error
}
