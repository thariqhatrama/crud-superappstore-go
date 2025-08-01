package service

import (
	"context"
	"errors"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

// UpdateUserRequest represents payload for updating user profile
type UpdateUserRequest struct {
	Nama         string
	TanggalLahir string
	JenisKelamin string
	Tentang      string
	Pekerjaan    string
	IDProvinsi   string
	IDKota       string
}

// UserService defines methods for user profile management
type UserService interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, id uint, req UpdateUserRequest) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService constructs a UserService with injected UserRepository
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, req UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Apply updates
	user.Nama = req.Nama
	// user.TanggalLahir = req.TanggalLahir
	user.JenisKelamin = req.JenisKelamin
	user.Tentang = req.Tentang
	user.Pekerjaan = req.Pekerjaan
	user.IDProvinsi = req.IDProvinsi
	user.IDKota = req.IDKota

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
