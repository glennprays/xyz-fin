package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/stretchr/testify/suite"
)

type transactionRepositoryTestSuite struct {
	suite.Suite
	DB   *sql.DB
	Mock sqlmock.Sqlmock
	Repo TransactionRepository
}

func (s *transactionRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)

	s.DB = db
	s.Mock = mock
	s.Repo = NewTransactionRepository(db)
}

func (s *transactionRepositoryTestSuite) TearDownTest() {
	s.DB.Close()
}

func (s *transactionRepositoryTestSuite) TestSave_Success() {
	ctx := context.Background()
	s.Mock.ExpectBegin()
	tx, err := s.DB.Begin()
	s.Require().NoError(err)

	transaction := &model.Transaction{
		NomorKontrak:  "TRX12345",
		ConsumerNIK:   "1234567890",
		OTR:           100000000,
		AdminFee:      500000,
		JumlahCicilan: 24,
		JumlahBunga:   5.5,
		NamaAsset:     "Motor Beat",
		Status:        "pending",
	}

	query := `INSERT INTO transactions \(nomor_kontrak, consumer_nik, otr, admin_fee, jumlah_cicilan, jumlah_bunga, nama_asset, status\) 
			  VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`

	s.Mock.ExpectExec(query).
		WithArgs(
			transaction.NomorKontrak,
			transaction.ConsumerNIK,
			transaction.OTR,
			transaction.AdminFee,
			transaction.JumlahCicilan,
			transaction.JumlahBunga,
			transaction.NamaAsset,
			transaction.Status,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.Repo.Save(ctx, tx, transaction)
	s.Require().NoError(err)

	s.Mock.ExpectCommit()
	err = tx.Commit()
	s.Require().NoError(err)
}

func (s *transactionRepositoryTestSuite) TestSave_Error() {
	ctx := context.Background()
	s.Mock.ExpectBegin()
	tx, err := s.DB.Begin()
	s.Require().NoError(err)

	transaction := &model.Transaction{
		NomorKontrak:  "TRX54321",
		ConsumerNIK:   "0987654321",
		OTR:           200000000,
		AdminFee:      1000000,
		JumlahCicilan: 36,
		JumlahBunga:   6.5,
		NamaAsset:     "Mobil Avanza",
		Status:        "pending",
	}

	query := `INSERT INTO transactions \(nomor_kontrak, consumer_nik, otr, admin_fee, jumlah_cicilan, jumlah_bunga, nama_asset, status\) 
			  VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`

	s.Mock.ExpectExec(query).
		WithArgs(
			transaction.NomorKontrak,
			transaction.ConsumerNIK,
			transaction.OTR,
			transaction.AdminFee,
			transaction.JumlahCicilan,
			transaction.JumlahBunga,
			transaction.NamaAsset,
			transaction.Status,
		).
		WillReturnError(sql.ErrConnDone)

	err = s.Repo.Save(ctx, tx, transaction)
	s.Require().Error(err)
	s.Equal(sql.ErrConnDone, err)

	s.Mock.ExpectRollback()
	err = tx.Rollback()
	s.Require().NoError(err)
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(transactionRepositoryTestSuite))
}
