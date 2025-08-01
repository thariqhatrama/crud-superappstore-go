package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"FinalTask/internal/models"
	"FinalTask/internal/repository"
)

type CreateTransactionRequest struct {
	AlamatPengiriman uint
	Items            []struct {
		LogProdukID uint
		Kuantitas   int
	}
	MethodBayar string
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
	// 1. Buat header transaksi
	trx := &models.Trx{
		IDUser:           userID,
		AlamatPengiriman: req.AlamatPengiriman,
		MethodBayar:      req.MethodBayar,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := s.trxRepo.Create(ctx, trx); err != nil {
		return nil, err
	}

	total := 0
	// 2. Proses setiap item: baca log, simpan detail, update stok
	for _, item := range req.Items {
		logEntry, err := s.productRepo.FindLogByID(ctx, item.LogProdukID)
		if err != nil {
			return nil, err
		}
		price, err := strconv.Atoi(logEntry.HargaKonsumen)
		if err != nil {
			return nil, fmt.Errorf("invalid harga konsumen: %v", err)
		}

		detail := &models.DetailTrx{
			IDTrx:       trx.ID,
			IDLogProduk: logEntry.ID,
			IDToko:      logEntry.IDToko,
			Kuantitas:   item.Kuantitas,
			HargaTotal:  item.Kuantitas * price,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := s.trxRepo.CreateDetail(ctx, detail); err != nil {
			return nil, err
		}
		total += detail.HargaTotal

		// Kurangi stok produk
		if err := s.productRepo.UpdateStock(ctx, logEntry.IDProduk, item.Kuantitas); err != nil {
			return nil, err
		}
	}

	// 3. Update total & kode invoice
	trx.HargaTotal = total
	trx.KodeInvoice = fmt.Sprintf("INV-%d", trx.ID)
	if err := s.trxRepo.Update(ctx, trx); err != nil {
		return nil, err
	}

	return trx, nil
}

func (s *transactionService) List(ctx context.Context, userID uint, qs map[string]string) ([]*models.Trx, error) {
	page, limit := 1, 10
	if p, ok := qs["page"]; ok {
		fmt.Sscanf(p, "%d", &page)
	}
	if l, ok := qs["limit"]; ok {
		fmt.Sscanf(l, "%d", &limit)
	}
	return s.trxRepo.ListByUserID(ctx, userID, (page-1)*limit, limit)
}

func (s *transactionService) GetByID(ctx context.Context, userID, id uint) (*models.Trx, error) {
	trx, err := s.trxRepo.FindByID(ctx, userID, id)
	if err != nil {
		return nil, errors.New("transaction not found or unauthorized")
	}
	return trx, nil
}
