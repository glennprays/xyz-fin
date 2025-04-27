package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type consumerLimitRepositoryTestSuite struct {
	suite.Suite
	DB   *sql.DB
	Mock sqlmock.Sqlmock
	Repo ConsumerLimitRepository
}

func (s *consumerLimitRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	s.Require().NoError(err)
	s.DB = db
	s.Mock = mock
	s.Repo = NewConsumerLimitRepository(db)
}

func (s *consumerLimitRepositoryTestSuite) TearDownTest() {
	s.DB.Close()
}

func (s *consumerLimitRepositoryTestSuite) TestFindByNIK() {
	nik := "12345"
	dummyTime := time.Now()

	rows := sqlmock.NewRows([]string{
		"consumer_nik", "tenor", "limit_amount", "created_at", "updated_at",
	}).AddRow(
		"12345", 6, 500000.00, dummyTime, dummyTime,
	).AddRow(
		"12345", 12, 1000000.00, dummyTime, dummyTime,
	).AddRow(
		"12345", 24, 2000000.00, dummyTime, dummyTime,
	)

	s.Mock.ExpectQuery(`SELECT consumer_nik, tenor, limit_amount, created_at, updated_at FROM consumer_limits WHERE nik = \$1 ORDER BY tenor ASC`).
		WithArgs(nik).
		WillReturnRows(rows)

	ctx := context.Background()

	consumerLimits, err := s.Repo.FindByNIK(ctx, nik)

	s.Require().NoError(err)
	s.Require().NotNil(consumerLimits)
	s.Len(consumerLimits, 3)
	s.Equal("12345", consumerLimits[0].ConsumerNIK)
	s.Equal(6, consumerLimits[0].Tenor)
	s.Equal(500000.00, consumerLimits[0].LimitAmount)
	s.Equal("12345", consumerLimits[1].ConsumerNIK)
	s.Equal(12, consumerLimits[1].Tenor)
	s.Equal(1000000.00, consumerLimits[1].LimitAmount)
	s.Equal("12345", consumerLimits[2].ConsumerNIK)
	s.Equal(24, consumerLimits[2].Tenor)
	s.Equal(2000000.00, consumerLimits[2].LimitAmount)
}

func (s *consumerLimitRepositoryTestSuite) TestFindByNIKAndTenor_Success() {
	nik := "12345"
	tenor := 12
	dummyTime := time.Now()

	s.Mock.ExpectBegin()
	tx, err := s.DB.Begin()
	s.Require().NoError(err)

	row := sqlmock.NewRows([]string{
		"consumer_nik", "tenor", "limit_amount", "created_at", "updated_at",
	}).AddRow(
		nik, tenor, 1000000.00, dummyTime, dummyTime,
	)

	s.Mock.ExpectQuery(`SELECT consumer_nik, tenor, limit_amount, created_at, updated_at FROM consumer_limits WHERE consumer_nik = \$1 AND tenor = \$2`).
		WithArgs(nik, tenor).
		WillReturnRows(row)

	ctx := context.Background()

	consumerLimit, err := s.Repo.FindByNIKAndTenor(ctx, tx, nik, tenor)
	s.Require().NoError(err)
	s.Require().NotNil(consumerLimit)
	s.Equal(nik, consumerLimit.ConsumerNIK)
	s.Equal(tenor, consumerLimit.Tenor)
	s.Equal(1000000.00, consumerLimit.LimitAmount)

	s.Mock.ExpectCommit()
	err = tx.Commit()
	s.Require().NoError(err)
}

func (s *consumerLimitRepositoryTestSuite) TestFindByNIKAndTenor_NotFound() {
	nik := "notfoundnik"
	tenor := 99

	s.Mock.ExpectBegin()
	tx, err := s.DB.Begin()
	s.Require().NoError(err)

	s.Mock.ExpectQuery(`SELECT consumer_nik, tenor, limit_amount, created_at, updated_at FROM consumer_limits WHERE consumer_nik = \$1 AND tenor = \$2`).
		WithArgs(nik, tenor).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()

	consumerLimit, err := s.Repo.FindByNIKAndTenor(ctx, tx, nik, tenor)
	s.Require().NoError(err)
	s.Require().Nil(consumerLimit)

	s.Mock.ExpectCommit()
	err = tx.Commit()
	s.Require().NoError(err)
}

func TestConsumerLimitRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(consumerLimitRepositoryTestSuite))
}
