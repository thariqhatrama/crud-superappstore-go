package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

// ==== Request DTO ====
type CreateAddressRequest struct {
	JudulAlamat  string `json:"judul_alamat"`
	NamaPenerima string `json:"nama_penerima"`
	NoTelp       string `json:"no_telp"`
	DetailAlamat string `json:"detail_alamat"`
}

// ==== Interface ====
type AddressService interface {
	Create(ctx context.Context, userID uint, req CreateAddressRequest) (*models.Alamat, error)
	List(ctx context.Context, userID uint) ([]*models.Alamat, error)
	Update(ctx context.Context, userID, id uint, req CreateAddressRequest) (*models.Alamat, error)
	Delete(ctx context.Context, userID, id uint) error
	GetByID(ctx context.Context, userID, id uint) (*models.Alamat, error)

	// Tambahan untuk API Emsifa
	GetProvinces(ctx context.Context) (interface{}, error)
	GetRegenciesByProvince(ctx context.Context, provinceID string) (interface{}, error)
}

// ==== Implementasi ====
type addressService struct {
	repo repository.AddressRepository
}

func NewAddressService(repo repository.AddressRepository) AddressService {
	return &addressService{repo: repo}
}

// ==== CRUD Alamat ====
func (s *addressService) Create(ctx context.Context, userID uint, req CreateAddressRequest) (*models.Alamat, error) {
	addr := &models.Alamat{
		IDUser:       userID,
		JudulAlamat:  req.JudulAlamat,
		NamaPenerima: req.NamaPenerima,
		NoTelp:       req.NoTelp,
		DetailAlamat: req.DetailAlamat,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repo.Create(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

func (s *addressService) List(ctx context.Context, userID uint) ([]*models.Alamat, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *addressService) Update(ctx context.Context, userID, id uint, req CreateAddressRequest) (*models.Alamat, error) {
	addr, err := s.repo.FindByID(ctx, id)
	if err != nil || addr.IDUser != userID {
		return nil, errors.New("address not found or unauthorized")
	}
	addr.JudulAlamat = req.JudulAlamat
	addr.NamaPenerima = req.NamaPenerima
	addr.NoTelp = req.NoTelp
	addr.DetailAlamat = req.DetailAlamat
	addr.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, addr); err != nil {
		return nil, err
	}
	return addr, nil
}

func (s *addressService) Delete(ctx context.Context, userID, id uint) error {
	addr, err := s.repo.FindByID(ctx, id)
	if err != nil || addr.IDUser != userID {
		return errors.New("address not found or unauthorized")
	}
	return s.repo.Delete(ctx, id)
}

func (s *addressService) GetByID(ctx context.Context, userID, id uint) (*models.Alamat, error) {
	addr, err := s.repo.FindByID(ctx, id)
	if err != nil || addr.IDUser != userID {
		return nil, errors.New("address not found or unauthorized")
	}
	return addr, nil
}

// ==== API Wilayah Indonesia ====
func (s *addressService) GetProvinces(ctx context.Context) (interface{}, error) {
	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var provinces interface{}
	if err := json.NewDecoder(resp.Body).Decode(&provinces); err != nil {
		return nil, err
	}
	return provinces, nil
}

func (s *addressService) GetRegenciesByProvince(ctx context.Context, provinceID string) (interface{}, error) {
	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/" + provinceID + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var regencies interface{}
	if err := json.NewDecoder(resp.Body).Decode(&regencies); err != nil {
		return nil, err
	}
	return regencies, nil
}
