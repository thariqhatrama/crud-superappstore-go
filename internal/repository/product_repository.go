package repository

import (
	"context"

	"FinalTask/config"
	"FinalTask/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, prod *models.Produk) error
	List(ctx context.Context, offset, limit int, categoryID uint) ([]*models.Produk, error)
	FindByID(ctx context.Context, id uint) (*models.Produk, error)
	Update(ctx context.Context, prod *models.Produk) error
	Delete(ctx context.Context, id uint) error
	CreatePhoto(ctx context.Context, photo *models.FotoProduk) error

	// Added for transactions
	FindLogByID(ctx context.Context, logID uint) (*models.LogProduk, error)
	UpdateStock(ctx context.Context, produkID uint, qty int) error
}

type productRepo struct{}

func NewProductRepository() ProductRepository {
	return &productRepo{}
}

func (r *productRepo) Create(ctx context.Context, prod *models.Produk) error {
	return config.DB.WithContext(ctx).Create(prod).Error
}

func (r *productRepo) List(ctx context.Context, offset, limit int, categoryID uint) ([]*models.Produk, error) {
	var list []*models.Produk
	db := config.DB.WithContext(ctx)
	if categoryID != 0 {
		db = db.Where("id_category = ?", categoryID)
	}
	err := db.Offset(offset).
		Limit(limit).
		Preload("FotoProduk").
		Find(&list).Error
	return list, err
}

func (r *productRepo) FindByID(ctx context.Context, id uint) (*models.Produk, error) {
	var prod models.Produk
	err := config.DB.WithContext(ctx).
		Preload("FotoProduk").
		First(&prod, id).Error
	return &prod, err
}

func (r *productRepo) Update(ctx context.Context, prod *models.Produk) error {
	return config.DB.WithContext(ctx).Save(prod).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint) error {
	return config.DB.WithContext(ctx).Delete(&models.Produk{}, id).Error
}

func (r *productRepo) CreatePhoto(ctx context.Context, photo *models.FotoProduk) error {
	return config.DB.WithContext(ctx).Create(photo).Error
}

// FindLogByID retrieves a LogProduk entry by its ID
func (r *productRepo) FindLogByID(ctx context.Context, logID uint) (*models.LogProduk, error) {
	var logEntry models.LogProduk
	if err := config.DB.WithContext(ctx).First(&logEntry, logID).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// UpdateStock decreases the 'stok' field of a Produk by the given qty
func (r *productRepo) UpdateStock(ctx context.Context, produkID uint, qty int) error {
	return config.DB.WithContext(ctx).
		Model(&models.Produk{}).
		Where("id = ?", produkID).
		UpdateColumn("stok", gorm.Expr("stok - ?", qty)).Error
}
