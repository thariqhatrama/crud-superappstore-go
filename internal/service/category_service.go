package service

import (
	"context"
	"errors"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

type CreateCategoryRequest struct {
	NamaCategory string `json:"nama_category"`
}

type UpdateCategoryRequest struct {
	NamaCategory string `json:"nama_category"`
}

type CategoryService interface {
	Create(ctx context.Context, req CreateCategoryRequest) (*models.Category, error)
	List(ctx context.Context) ([]*models.Category, error)
	Update(ctx context.Context, id uint, req UpdateCategoryRequest) (*models.Category, error)
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*models.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (s *categoryService) Create(ctx context.Context, req CreateCategoryRequest) (*models.Category, error) {
	cat := &models.Category{NamaCategory: req.NamaCategory}
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) List(ctx context.Context) ([]*models.Category, error) {
	return s.repo.List(ctx)
}

func (s *categoryService) Update(ctx context.Context, id uint, req UpdateCategoryRequest) (*models.Category, error) {
	cat, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	cat.NamaCategory = req.NamaCategory
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return errors.New("category not found")
	}
	return s.repo.Delete(ctx, id)
}
