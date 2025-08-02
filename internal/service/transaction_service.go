package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"FinalTask/config"
	"FinalTask/internal/models"
	"FinalTask/internal/repository"

	"gorm.io/gorm"
)

type CreateTransactionRequest struct {
	AlamatPengiriman uint `json:"alamat_pengiriman"`
	Items            []struct {
		LogProdukID uint `json:"log_produk_id"`
		Kuantitas   int  `json:"kuantitas"`
	} `json:"items"`
	MethodBayar string `json:"method_bayar"`
}

type TransactionService interface {
	Create(ctx context.Context, userID uint, req CreateTransactionRequest) (*models.Trx, error)
	List(ctx context.Context, userID uint, qs map[string]string) ([]*models.Trx, error)
	GetByID(ctx context.Context, userID, id uint) (*models.Trx, error)
}

type transactionService struct {
	trxRepo     repository.TransactionRepository
	productRepo repository.ProductRepository
}

func NewTransactionService(
	trxRepo repository.TransactionRepository,
	productRepo repository.ProductRepository,
) TransactionService {
	return &transactionService{
		trxRepo:     trxRepo,
		productRepo: productRepo,
	}
}

func (s *transactionService) Create(ctx context.Context, userID uint, req CreateTransactionRequest) (*models.Trx, error) {
	var result *models.Trx

	err := config.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Generate kode invoice unik
		invoiceCode := fmt.Sprintf("INV-%d-%d", time.Now().UnixNano(), userID)

		now := time.Now()
		trx := &models.Trx{
			IDUser:           userID,
			AlamatPengiriman: req.AlamatPengiriman,
			MethodBayar:      req.MethodBayar,
			KodeInvoice:      invoiceCode,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		// simpan header transaksi
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		total := 0
		// 2. Proses tiap item
		for _, item := range req.Items {
			// baca snapshot log produk
			logEntry, err := s.productRepo.FindLogByID(ctx, item.LogProdukID)
			if err != nil {
				return err
			}
			price, err := strconv.Atoi(logEntry.HargaKonsumen)
			if err != nil {
				return fmt.Errorf("invalid harga konsumen: %v", err)
			}

			// buat detail transaksi
			detail := &models.DetailTrx{
				IDTrx:       trx.ID,
				IDLogProduk: logEntry.ID,
				IDToko:      logEntry.IDToko,
				Kuantitas:   item.Kuantitas,
				HargaTotal:  item.Kuantitas * price,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := tx.Create(detail).Error; err != nil {
				return err
			}
			total += detail.HargaTotal

			// kurangi stok di produk
			if err := s.productRepo.UpdateStock(ctx, logEntry.IDProduk, item.Kuantitas); err != nil {
				return err
			}
		}

		// 3. Update total harga dan timestamp
		trx.HargaTotal = total
		trx.UpdatedAt = time.Now()
		if err := tx.Save(trx).Error; err != nil {
			return err
		}

		result = trx
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *transactionService) List(ctx context.Context, userID uint, qs map[string]string) ([]*models.Trx, error) {
	page, limit := 1, 10
	if p, ok := qs["page"]; ok {
		fmt.Sscanf(p, "%d", &page)
	}
	if l, ok := qs["limit"]; ok {
		fmt.Sscanf(l, "%d", &limit)
	}
	// ListByUserID should preload DetailTrx and related
	return s.trxRepo.ListByUserID(ctx, userID, (page-1)*limit, limit)
}

func (s *transactionService) GetByID(ctx context.Context, userID, id uint) (*models.Trx, error) {
	trx, err := s.trxRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, errors.New("transaction not found or unauthorized")
	}
	return trx, nil
}
