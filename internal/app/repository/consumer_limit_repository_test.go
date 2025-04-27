package repository

import (
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
		"12345", 12, 1000000.00, dummyTime, dummyTime,
	)

	s.Mock.ExpectQuery(`SELECT consumer_nik, tenor, limit_amount, created_at, updated_at FROM consumer_limits WHERE nik = \$1`).
		WithArgs(nik).
		WillReturnRows(rows)

	consumerLimit, err := s.Repo.FindByNIK(nik)

	s.Require().NoError(err)
	s.Require().NotNil(consumerLimit)
	s.Equal("12345", consumerLimit.ConsumerNIK)
	s.Equal(12, consumerLimit.Tenor)
	s.Equal(1000000.00, consumerLimit.LimitAmount)
}

func TestConsumerLimitRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(consumerLimitRepositoryTestSuite))
}
