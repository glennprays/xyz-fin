package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/glennprays/xyz-fin/internal/app/repository"
	"github.com/glennprays/xyz-fin/internal/app/service"
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, nik string, req *model.TransactionRequest) (*model.TransactionResponse, error)
}

type transactionUsecase struct {
	db                 *sql.DB
	transactionService service.TransactionService
	transactionRepo    repository.TransactionRepository
	consumerRepo       repository.ConsumerRepository
	consumerLimitRepo  repository.ConsumerLimitRepository
}

func NewTransactionUsecase(
	db *sql.DB,
	transactionService service.TransactionService,
	transactionRepo repository.TransactionRepository,
	consumerRepo repository.ConsumerRepository,
	consumerLimitRepo repository.ConsumerLimitRepository,
) TransactionUsecase {
	return &transactionUsecase{
		db:                 db,
		transactionService: transactionService,
		transactionRepo:    transactionRepo,
		consumerRepo:       consumerRepo,
		consumerLimitRepo:  consumerLimitRepo,
	}
}

func (u *transactionUsecase) CreateTransaction(ctx context.Context, nik string, req *model.TransactionRequest) (*model.TransactionResponse, error) {
	if nik != req.ConsumerNIK {
		appErr := errors.New("nik does not match with consumer NIK")
		return nil, model.NewError(model.ErrBadRequest, appErr)
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		appErr := errors.New("failed to begin transaction")
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	consumer, err := u.consumerRepo.FindAndLockByNIK(ctx, tx, req.ConsumerNIK)
	if err != nil {
		appErr := errors.New("failed to find consumer")
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}
	if consumer == nil {
		appErr := errors.New("consumer not found")
		return nil, model.NewError(model.ErrNotFound, appErr)
	}

	consumerLimits, err := u.consumerLimitRepo.FindByNIKAndTenor(ctx, tx, nik, req.Tenor)
	if err != nil {
		appErr := errors.New("failed to find consumer limits")
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}
	if consumerLimits == nil {
		appErr := errors.New("consumer limits not found")
		return nil, model.NewError(model.ErrNotFound, appErr)
	}

	activeSUM, err := u.transactionRepo.GetActiveTransactionSumByNIK(ctx, tx, nik)
	if err != nil {
		appErr := errors.New("failed to get active transaction sum")
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}

	if activeSUM+req.OTR > consumerLimits.LimitAmount {
		appErr := errors.New("transaction exceeds limit")
		return nil, model.NewError(model.ErrBadRequest, appErr)
	}

	transactionID := u.transactionService.GenerateTransactionID()
	// admin fee 3% and bungan 5%
	adminFee := req.OTR * 0.03
	jumlahBunga := req.OTR * 0.05
	transaction := &model.Transaction{
		NomorKontrak:  transactionID,
		ConsumerNIK:   req.ConsumerNIK,
		OTR:           req.OTR,
		AdminFee:      adminFee,
		JumlahBunga:   jumlahBunga,
		JumlahCicilan: req.Tenor,
		NamaAsset:     req.NamaAsset,
		Status:        "ACTIVE",
	}
	err = u.transactionRepo.Save(ctx, tx, transaction)
	if err != nil {
		appErr := errors.New("failed to save transaction")
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}

	transactionResponse := &model.TransactionResponse{
		NomorKontrak:  transaction.NomorKontrak,
		ConsumerNIK:   transaction.ConsumerNIK,
		OTR:           transaction.OTR,
		AdminFee:      transaction.AdminFee,
		JumlahCicilan: transaction.JumlahCicilan,
		JumlahBunga:   transaction.JumlahBunga,
		NamaAsset:     transaction.NamaAsset,
		Status:        transaction.Status,
	}

	return transactionResponse, nil
}
