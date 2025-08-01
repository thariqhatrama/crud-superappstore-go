package service

import (
	"context"
	"errors"
	"time"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
	"FinalTask/utils"

	"gorm.io/gorm"
)

type RegisterRequest struct {
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	NoTelp       string `json:"no_telp"`
	Password     string `json:"password"`
	TanggalLahir string `json:"tanggal_lahir"` // format: YYYY-MM-DD
	JenisKelamin string `json:"jenis_kelamin"`
	Tentang      string `json:"tentang"`
	Pekerjaan    string `json:"pekerjaan"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	storeRepo repository.StoreRepository
}

func NewAuthService(userRepo repository.UserRepository, storeRepo repository.StoreRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// ====== Cek email unik ======
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // error DB lain
	}
	if err == nil && existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// ====== Cek no_telp unik ======
	existingByPhone, err := s.userRepo.FindByPhone(ctx, req.NoTelp)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // error DB lain
	}
	if err == nil && existingByPhone != nil {
		return nil, errors.New("nomor telepon sudah terdaftar")
	}

	// ====== Hash password ======
	hashed, _ := utils.HashPassword(req.Password)

	// ====== Parsing tanggal lahir ======
	var dob *time.Time
	if req.TanggalLahir != "" {
		parsed, err := time.Parse("2006-01-02", req.TanggalLahir)
		if err == nil {
			dob = &parsed
		}
	}

	// ====== Buat user baru ======
	user := &models.User{
		Nama:         req.Nama,
		Email:        req.Email,
		NoTelp:       req.NoTelp,
		Password:     hashed,
		TanggalLahir: dob,
		JenisKelamin: req.JenisKelamin,
		Tentang:      req.Tentang,
		Pekerjaan:    req.Pekerjaan,
		IDProvinsi:   req.IDProvinsi,
		IDKota:       req.IDKota,
		IsAdmin:      false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// ====== Auto-create store ======
	store := &models.Toko{
		IDUser:    user.ID,
		NamaToko:  req.Nama + "'s Store",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.storeRepo.Create(ctx, store); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (string, error) {
	// Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// Cek password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("email atau password salah")
	}

	// Buat JWT
	token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
	return token, err
}
