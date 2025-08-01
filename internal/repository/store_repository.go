package repository

import (
	"FinalTask/config"
	"FinalTask/internal/models"
	"context"
)

type StoreRepository interface {
	FindByUserID(ctx context.Context, userID uint) (*models.Toko, error)
	Create(ctx context.Context, store *models.Toko) error
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

func (r *storeRepo) Create(ctx context.Context, store *models.Toko) error {
	return config.DB.WithContext(ctx).Create(store).Error
}

func (r *storeRepo) Update(ctx context.Context, store *models.Toko) error {
	return config.DB.WithContext(ctx).Save(store).Error
}
