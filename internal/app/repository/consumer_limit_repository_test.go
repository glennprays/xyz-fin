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

func TestConsumerLimitRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(consumerLimitRepositoryTestSuite))
}
