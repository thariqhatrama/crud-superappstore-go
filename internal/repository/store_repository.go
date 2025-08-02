package repository

import (
	"context"

	"FinalTask/config"
	"FinalTask/internal/models"
)

type StoreRepository interface {
	// Find store by owner user ID
	FindByUserID(ctx context.Context, userID uint) (*models.Toko, error)
	// Find store by its own ID
	FindByID(ctx context.Context, id uint) (*models.Toko, error)
	// List all stores (for admin)
	List(ctx context.Context) ([]*models.Toko, error)
	// Create new store
	Create(ctx context.Context, store *models.Toko) error
	// Update existing store
	Update(ctx context.Context, store *models.Toko) error
}

type storeRepo struct{}

func NewStoreRepository() StoreRepository {
	return &storeRepo{}
}

func (r *storeRepo) FindByUserID(ctx context.Context, userID uint) (*models.Toko, error) {
	var store models.Toko
	err := config.DB.WithContext(ctx).
		Where("id_user = ?", userID).
		First(&store).Error
	return &store, err
}

func (r *storeRepo) FindByID(ctx context.Context, id uint) (*models.Toko, error) {
	var store models.Toko
	err := config.DB.WithContext(ctx).
		First(&store, id).Error
	return &store, err
}

func (r *storeRepo) List(ctx context.Context) ([]*models.Toko, error) {
	var stores []*models.Toko
	err := config.DB.WithContext(ctx).
		Find(&stores).Error
	return stores, err
}

func (r *storeRepo) Create(ctx context.Context, store *models.Toko) error {
	return config.DB.WithContext(ctx).Create(store).Error
}

func (r *storeRepo) Update(ctx context.Context, store *models.Toko) error {
	return config.DB.WithContext(ctx).Save(store).Error
}
