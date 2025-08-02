package service

import (
	"context"
	"errors"
	"time"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

type UpdateStoreRequest struct {
	NamaToko string `json:"nama_toko"`
	URLFoto  string `json:"url_foto"`
}

type StoreService interface {
	// Untuk endpoint GET  /store
	GetByUser(ctx context.Context, userID uint) (*models.Toko, error)
	// Untuk endpoint PUT  /store
	Update(ctx context.Context, userID uint, req UpdateStoreRequest) (*models.Toko, error)
	// Untuk endpoint GET  /stores
	ListAll(ctx context.Context) ([]*models.Toko, error)
	// Untuk endpoint GET  /stores/:id
	GetByID(ctx context.Context, id uint) (*models.Toko, error)
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{repo: repo}
}

func (s *storeService) GetByUser(ctx context.Context, userID uint) (*models.Toko, error) {
	store, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	return store, nil
}

func (s *storeService) Update(ctx context.Context, userID uint, req UpdateStoreRequest) (*models.Toko, error) {
	store, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	store.NamaToko = req.NamaToko
	store.URLFoto = req.URLFoto
	store.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *storeService) ListAll(ctx context.Context) ([]*models.Toko, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve stores")
	}
	return list, nil
}

func (s *storeService) GetByID(ctx context.Context, id uint) (*models.Toko, error) {
	store, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("store not found")
	}
	return store, nil
}
