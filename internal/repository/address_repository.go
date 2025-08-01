package repository

import (
	"FinalTask/config"
	"FinalTask/internal/models"
	"context"
)

type AddressRepository interface {
	Create(ctx context.Context, addr *models.Alamat) error
	ListByUserID(ctx context.Context, userID uint) ([]*models.Alamat, error)
	FindByID(ctx context.Context, id uint) (*models.Alamat, error)
	Update(ctx context.Context, addr *models.Alamat) error
	Delete(ctx context.Context, id uint) error
}

type addressRepo struct{}

func NewAddressRepository() AddressRepository {
	return &addressRepo{}
}

func (r *addressRepo) Create(ctx context.Context, addr *models.Alamat) error {
	return config.DB.WithContext(ctx).Create(addr).Error
}

func (r *addressRepo) ListByUserID(ctx context.Context, userID uint) ([]*models.Alamat, error) {
	var list []*models.Alamat
	err := config.DB.WithContext(ctx).
		Where("id_user = ?", userID).
		Find(&list).Error
	return list, err
}

func (r *addressRepo) FindByID(ctx context.Context, id uint) (*models.Alamat, error) {
	var addr models.Alamat
	err := config.DB.WithContext(ctx).First(&addr, id).Error
	return &addr, err
}

func (r *addressRepo) Update(ctx context.Context, addr *models.Alamat) error {
	return config.DB.WithContext(ctx).Save(addr).Error
}

func (r *addressRepo) Delete(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Delete(&models.Alamat{}, id).Error
}
