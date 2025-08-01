package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

type CreateProductRequest struct {
	NamaProduk    string `json:"nama_produk"`
	Slug          string `json:"slug"`
	HargaReseller string `json:"harga_reseller"`
	HargaKonsumen string `json:"harga_konsumen"`
	Stok          int    `json:"stok"`
	Deskripsi     string `json:"deskripsi"`
	IDCategory    uint   `json:"id_category"`
}

type ProductService interface {
	Create(ctx context.Context, userID uint, req CreateProductRequest) (*models.Produk, error)
	List(ctx context.Context, qs map[string]string) ([]*models.Produk, error)
	GetByID(ctx context.Context, id uint) (*models.Produk, error)
	Update(ctx context.Context, userID, id uint, req CreateProductRequest) (*models.Produk, error)
	Delete(ctx context.Context, userID, id uint) error
	UploadImage(ctx context.Context, id uint, file *multipart.FileHeader) (string, error)
}

type productService struct {
	repo         repository.ProductRepository
	storeRepo    repository.StoreRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(
	pr repository.ProductRepository,
	sr repository.StoreRepository,
	cr repository.CategoryRepository,
) ProductService {
	return &productService{
		repo:         pr,
		storeRepo:    sr,
		categoryRepo: cr,
	}
}

func (s *productService) Create(ctx context.Context, userID uint, req CreateProductRequest) (*models.Produk, error) {
	// Validasi toko user
	store, err := s.storeRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("store not found for user")
	}
	// Validasi kategori
	if _, err := s.categoryRepo.FindByID(ctx, req.IDCategory); err != nil {
		return nil, errors.New("category not found")
	}
	// Generate slug jika kosong
	slug := req.Slug
	if slug == "" {
		slug = fmt.Sprintf("%s-%d", req.NamaProduk, time.Now().Unix())
	}

	prod := &models.Produk{
		NamaProduk:    req.NamaProduk,
		Slug:          slug,
		HargaReseller: req.HargaReseller,
		HargaKonsumen: req.HargaKonsumen,
		Stok:          req.Stok,
		Deskripsi:     req.Deskripsi,
		IDToko:        store.ID,
		IDCategory:    req.IDCategory,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, prod); err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *productService) List(ctx context.Context, qs map[string]string) ([]*models.Produk, error) {
	page := 1
	limit := 10
	category := uint(0)

	if v, ok := qs["page"]; ok {
		if p, err := strconv.Atoi(v); err == nil {
			page = p
		}
	}
	if v, ok := qs["limit"]; ok {
		if l, err := strconv.Atoi(v); err == nil {
			limit = l
		}
	}
	if v, ok := qs["id_kategori"]; ok {
		if c, err := strconv.Atoi(v); err == nil {
			category = uint(c)
		}
	}

	return s.repo.List(ctx, (page-1)*limit, limit, category)
}

func (s *productService) GetByID(ctx context.Context, id uint) (*models.Produk, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) Update(ctx context.Context, userID, id uint, req CreateProductRequest) (*models.Produk, error) {
	// Ambil produk dan periksa kepemilikan
	prod, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	store, err := s.storeRepo.FindByUserID(ctx, userID)
	if err != nil || prod.IDToko != store.ID {
		return nil, errors.New("unauthorized")
	}
	// Validasi kategori
	if _, err := s.categoryRepo.FindByID(ctx, req.IDCategory); err != nil {
		return nil, errors.New("category not found")
	}

	prod.NamaProduk = req.NamaProduk
	if req.Slug != "" {
		prod.Slug = req.Slug
	}
	prod.HargaReseller = req.HargaReseller
	prod.HargaKonsumen = req.HargaKonsumen
	prod.Stok = req.Stok
	prod.Deskripsi = req.Deskripsi
	prod.IDCategory = req.IDCategory
	prod.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, prod); err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *productService) Delete(ctx context.Context, userID, id uint) error {
	prod, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.New("product not found")
	}
	store, err := s.storeRepo.FindByUserID(ctx, userID)
	if err != nil || prod.IDToko != store.ID {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(ctx, id)
}

func (s *productService) UploadImage(ctx context.Context, id uint, file *multipart.FileHeader) (string, error) {
	// 1. Simpan file
	filename := fmt.Sprintf("%d_%s", id, filepath.Base(file.Filename))
	dest := filepath.Join("uploads", "products", filename)
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	// 2. Simpan record foto via ProductRepository
	photo := &models.FotoProduk{
		IDProduk:  id,
		URL:       "/" + filepath.ToSlash(dest),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.CreatePhoto(ctx, photo); err != nil {
		return "", err
	}

	return photo.URL, nil
}
