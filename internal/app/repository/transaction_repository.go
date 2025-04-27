package repository

import (
	"context"
	"database/sql"

	"github.com/glennprays/xyz-fin/internal/app/model"
)

type TransactionRepository interface {
	Save(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Save(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions (nomor_kontrak, consumer_nik, otr, admin_fee, jumlah_cicilan, jumlah_bunga, nama_asset, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.ExecContext(ctx, query,
		transaction.NomorKontrak,
		transaction.ConsumerNIK,
		transaction.OTR,
		transaction.AdminFee,
		transaction.JumlahCicilan,
		transaction.JumlahBunga,
		transaction.NamaAsset,
		transaction.Status,
	)
	return err
}
