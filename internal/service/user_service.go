package service

import (
	"context"
	"errors"
	"time"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

// UpdateUserRequest represents payload for updating user profile
type UpdateUserRequest struct {
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"` // expects "YYYY-MM-DD"
	JenisKelamin string `json:"jenis_kelamin"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
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

	// Update basic fields
	user.Nama = req.Nama
	user.JenisKelamin = req.JenisKelamin
	user.Tentang = req.Tentang
	user.Pekerjaan = req.Pekerjaan
	user.IDProvinsi = req.IDProvinsi
	user.IDKota = req.IDKota

	// Parse and update TanggalLahir if non-empty
	if req.TanggalLahir != "" {
		parsed, perr := time.Parse("2006-01-02", req.TanggalLahir)
		if perr != nil {
			return nil, errors.New("invalid tanggal_lahir format, expected YYYY-MM-DD")
		}
		user.TanggalLahir = &parsed
	}

	// Update UpdatedAt timestamp
	user.UpdatedAt = time.Now()

	// Persist changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
