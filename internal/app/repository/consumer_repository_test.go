package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type consumerRepositoryTestSuite struct {
	suite.Suite
	DB   *sql.DB
	Mock sqlmock.Sqlmock
	Repo ConsumerRepository
}

func (s *consumerRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.DB = db
	s.Mock = mock
	s.Repo = NewConsumerRepository(db)
}

func (s *consumerRepositoryTestSuite) TearDownTest() {
	s.DB.Close()
}

func (s *consumerRepositoryTestSuite) TestFindByPhoneNumber() {
	phone := "081234567890"
	dummyTime := time.Now()

	rows := sqlmock.NewRows([]string{
		"nik", "phone_number", "password_hash", "full_name", "legal_name", "tempat_lahir", "tanggal_lahir", "gaji", "foto_ktp_path", "foto_selfie_path", "created_at", "updated_at",
	}).AddRow(
		"12345", "081234567890", "hashed-password", "John Doe", "Legal Name", "City", dummyTime, 5000000.00, "/path/to/ktp.jpg", "/path/to/selfie.jpg", dummyTime, dummyTime,
	)

	s.Mock.ExpectQuery(`SELECT nik, phone_number, password_hash, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp_path, foto_selfie_path, created_at, updated_at FROM consumers WHERE phone_number = \$1`).
		WithArgs(phone).
		WillReturnRows(rows)

	ctx := context.Background()

	consumer, err := s.Repo.FindByPhoneNumber(ctx, phone)

	s.Require().NoError(err)
	s.Require().NotNil(consumer)
	s.Equal("12345", consumer.NIK)
	s.Equal("081234567890", consumer.PhoneNumber)
	s.Equal("John Doe", consumer.FullName)
}

func (s *consumerRepositoryTestSuite) TestFindByNIK() {
	nik := "12345"
	dummyTime := time.Now()

	rows := sqlmock.NewRows([]string{
		"nik", "phone_number", "password_hash", "full_name", "legal_name", "tempat_lahir", "tanggal_lahir", "gaji", "foto_ktp_path", "foto_selfie_path", "created_at", "updated_at",
	}).AddRow(
		"12345", "081234567890", "hashed-password", "John Doe", "Legal Name", "City", dummyTime, 5000000.00, "/path/to/ktp.jpg", "/path/to/selfie.jpg", dummyTime, dummyTime,
	)

	s.Mock.ExpectQuery(`SELECT nik, phone_number, password_hash, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp_path, foto_selfie_path, created_at, updated_at FROM consumers WHERE nik = \$1`).
		WithArgs(nik).
		WillReturnRows(rows)

	ctx := context.Background()

	consumer, err := s.Repo.FindByNIK(ctx, nik)

	s.Require().NoError(err)
	s.Require().NotNil(consumer)
	s.Equal("12345", consumer.NIK)
	s.Equal("081234567890", consumer.PhoneNumber)
	s.Equal("John Doe", consumer.FullName)
}

func (s *consumerRepositoryTestSuite) TestFindAndLockByNIK_Success() {
	ctx := context.Background()
	nik := "1234567890"
	dummyTime := time.Now()

	s.Mock.ExpectBegin()
	tx, err := s.DB.Begin()
	s.Require().NoError(err)

	rows := sqlmock.NewRows([]string{
		"nik", "phone_number", "password_hash", "full_name", "legal_name", "tempat_lahir", "tanggal_lahir", "gaji", "foto_ktp_path", "foto_selfie_path", "created_at", "updated_at",
	}).AddRow(
		nik, "08123456789", "hashedpassword", "Full Name", "Legal Name", "Tempat Lahir", dummyTime, 5000000.00, "ktp/path", "selfie/path", dummyTime, dummyTime,
	)

	s.Mock.ExpectQuery(`SELECT nik, phone_number, password_hash, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp_path, foto_selfie_path, created_at, updated_at FROM consumers WHERE nik = \$1 FOR UPDATE`).
		WithArgs(nik).
		WillReturnRows(rows)

	consumer, err := s.Repo.FindAndLockByNIK(ctx, tx, nik)

	s.Require().NoError(err)
	s.Require().NotNil(consumer)
	s.Equal(nik, consumer.NIK)
	s.Equal("08123456789", consumer.PhoneNumber)
	s.Equal("hashedpassword", consumer.PasswordHash)
	s.Equal("Full Name", consumer.FullName)
	s.Equal("Legal Name", consumer.LegalName)
	s.Equal("Tempat Lahir", consumer.TempatLahir)
	s.WithinDuration(dummyTime, consumer.TanggalLahir, time.Second)
	s.Equal(5000000.00, consumer.Gaji)
	s.Equal("ktp/path", consumer.FotoKTPPath)
	s.Equal("selfie/path", consumer.FotoSelfiePath)
	s.WithinDuration(dummyTime, consumer.CreatedAt, time.Second)
	s.WithinDuration(dummyTime, consumer.UpdatedAt, time.Second)

	s.Mock.ExpectCommit()
	err = tx.Commit()
	s.Require().NoError(err)
}

func (s *consumerRepositoryTestSuite) TestFindAndLockByNIK_NotFound() {
	ctx := context.Background()
	nik := "notfoundnik"

	s.Mock.ExpectBegin()
	tx, err := s.DB.Begin()
	s.Require().NoError(err)

	s.Mock.ExpectQuery(`SELECT nik, phone_number, password_hash, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp_path, foto_selfie_path, created_at, updated_at FROM consumers WHERE nik = \$1 FOR UPDATE`).
		WithArgs(nik).
		WillReturnError(sql.ErrNoRows)

	consumer, err := s.Repo.FindAndLockByNIK(ctx, tx, nik)

	s.Require().NoError(err)
	s.Nil(consumer)

	s.Mock.ExpectCommit()
	err = tx.Commit()
	s.Require().NoError(err)
}

func TestConsumerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(consumerRepositoryTestSuite))
}
