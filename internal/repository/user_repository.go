package repository

import (
	"context"

	"FinalTask/config"
	"FinalTask/internal/models"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}

type userRepo struct{}

func NewUserRepository() UserRepository {
	return &userRepo{}
}

// FindByID mencari user berdasarkan ID
func (r *userRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := config.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := config.DB.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone mencari user berdasarkan nomor telepon
func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	if err := config.DB.WithContext(ctx).
		Where("no_telp = ?", phone).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create membuat user baru
func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	return config.DB.WithContext(ctx).Create(user).Error
}

// Update memperbarui data user
func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	return config.DB.WithContext(ctx).Save(user).Error
}
