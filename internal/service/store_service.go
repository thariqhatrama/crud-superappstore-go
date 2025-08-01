package service

import (
	"context"
	"errors"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

type UpdateStoreRequest struct {
	NamaToko string
	URLFoto  string
}

type StoreService interface {
	GetByUser(ctx context.Context, userID uint) (*models.Toko, error)
	Update(ctx context.Context, userID uint, req UpdateStoreRequest) (*models.Toko, error)
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) GetByUser(ctx context.Context, userID uint) (*models.Toko, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *storeService) Update(ctx context.Context, userID uint, req UpdateStoreRequest) (*models.Toko, error) {
	store, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	store.NamaToko = req.NamaToko
	store.URLFoto = req.URLFoto
	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}
