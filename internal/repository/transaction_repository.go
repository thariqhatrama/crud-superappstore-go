package repository

import (
	"context"

	"FinalTask/config"
	"FinalTask/internal/models"
)

// TransactionRepository defines methods for handling transactions and detail rows
type TransactionRepository interface {
	Create(ctx context.Context, trx *models.Trx) error
	ListByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.Trx, error)
	FindByID(ctx context.Context, userID, id uint) (*models.Trx, error)

	// Methods needed by service layer
	CreateDetail(ctx context.Context, detail *models.DetailTrx) error
	Update(ctx context.Context, trx *models.Trx) error
}

// transactionRepo is concrete implementation of TransactionRepository
type transactionRepo struct{}

// NewTransactionRepository constructs a TransactionRepository
func NewTransactionRepository() TransactionRepository {
	return &transactionRepo{}
}

// Create inserts a new Trx record
func (r *transactionRepo) Create(ctx context.Context, trx *models.Trx) error {
	return config.DB.WithContext(ctx).Create(trx).Error
}

// ListByUserID returns a paginated list of Trx for a user
func (r *transactionRepo) ListByUserID(ctx context.Context, userID uint, offset, limit int) ([]*models.Trx, error) {
	var list []*models.Trx
	err := config.DB.WithContext(ctx).
		Where("id_user = ?", userID).
		Offset(offset).
		Limit(limit).
		Preload("DetailTrx").
		Find(&list).Error
	return list, err
}

// FindByID retrieves a single Trx by userID and trx ID
func (r *transactionRepo) FindByID(ctx context.Context, userID, id uint) (*models.Trx, error) {
	var trx models.Trx
	err := config.DB.WithContext(ctx).
		Where("id = ? AND id_user = ?", id, userID).
		Preload("DetailTrx").
		First(&trx).Error
	return &trx, err
}

// CreateDetail inserts a new DetailTrx record (transaction detail)
func (r *transactionRepo) CreateDetail(ctx context.Context, detail *models.DetailTrx) error {
	return config.DB.WithContext(ctx).Create(detail).Error
}

// Update saves changes to an existing Trx record (e.g., update total & invoice)
func (r *transactionRepo) Update(ctx context.Context, trx *models.Trx) error {
	return config.DB.WithContext(ctx).Save(trx).Error
}
