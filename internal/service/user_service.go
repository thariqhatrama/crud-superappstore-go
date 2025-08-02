package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
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

// UpdateProfileResult wraps the updated user plus region names
type UpdateProfileResult struct {
	User         *models.User `json:"user"`
	ProvinceName string       `json:"province_name"`
	CityName     string       `json:"city_name"`
}

// UserService defines methods for user profile management
type UserService interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, id uint, req UpdateUserRequest) (*UpdateProfileResult, error)
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

func (s *userService) Update(ctx context.Context, id uint, req UpdateUserRequest) (*UpdateProfileResult, error) {
	// 1) Load existing user
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2) Fetch and validate provinces
	provinces, err := fetchRegions("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return nil, fmt.Errorf("cannot load provinces: %v", err)
	}
	var provinceName string
	for _, p := range provinces {
		if p.ID == req.IDProvinsi {
			provinceName = p.Name
			break
		}
	}
	if provinceName == "" {
		return nil, errors.New("province not found")
	}

	// 3) Fetch and validate regencies for that province
	regURL := fmt.Sprintf("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/%s.json", req.IDProvinsi)
	regencies, err := fetchRegions(regURL)
	if err != nil {
		return nil, fmt.Errorf("cannot load regencies: %v", err)
	}
	var cityName string
	for _, r := range regencies {
		if r.ID == req.IDKota {
			cityName = r.Name
			break
		}
	}
	if cityName == "" {
		return nil, errors.New("city not found in selected province")
	}

	// 4) Apply updates to user struct
	user.Nama = req.Nama
	user.JenisKelamin = req.JenisKelamin
	user.Tentang = req.Tentang
	user.Pekerjaan = req.Pekerjaan
	user.IDProvinsi = req.IDProvinsi
	user.IDKota = req.IDKota
	if req.TanggalLahir != "" {
		if dt, perr := time.Parse("2006-01-02", req.TanggalLahir); perr == nil {
			user.TanggalLahir = &dt
		} else {
			return nil, errors.New("invalid tanggal_lahir format, expected YYYY-MM-DD")
		}
	}
	user.UpdatedAt = time.Now()

	// 5) Persist changes
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	// 6) Return enriched result
	return &UpdateProfileResult{
		User:         user,
		ProvinceName: provinceName,
		CityName:     cityName,
	}, nil
}

// regionEntry is used to decode Emsifa region JSON (id as string)
type regionEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// custom HTTP client with reasonable timeouts
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 3 * time.Second,
	},
}

// fetchRegions fetches and decodes a list of regions from the given URL
func fetchRegions(url string) ([]regionEntry, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var list []regionEntry
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}
