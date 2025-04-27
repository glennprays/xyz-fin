package repository

import (
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

	consumer, err := s.Repo.FindByPhoneNumber(phone)

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

	consumer, err := s.Repo.FindByNIK(nik)

	s.Require().NoError(err)
	s.Require().NotNil(consumer)
	s.Equal("12345", consumer.NIK)
	s.Equal("081234567890", consumer.PhoneNumber)
	s.Equal("John Doe", consumer.FullName)
}

func TestConsumerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(consumerRepositoryTestSuite))
}
