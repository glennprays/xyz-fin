package repository

import (
	"database/sql"

	"github.com/glennprays/xyz-fin/internal/app/model"
)

type ConsumerLimitRepository interface {
	FindByNIK(nik string) (*model.ConsumerLimit, error)
}

type consumerLimitRepository struct {
	db *sql.DB
}

func NewConsumerLimitRepository(db *sql.DB) ConsumerLimitRepository {
	return &consumerLimitRepository{db: db}
}

func (r *consumerLimitRepository) FindByNIK(nik string) (*model.ConsumerLimit, error) {
	consumerLimit := &model.ConsumerLimit{}
	query := `SELECT consumer_nik, tenor, limit_amount, created_at, updated_at
  FROM consumer_limits WHERE nik = $1`
	err := r.db.QueryRow(query, nik).Scan(&consumerLimit.ConsumerNIK, &consumerLimit.Tenor, &consumerLimit.LimitAmount, &consumerLimit.CreatedAt, &consumerLimit.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return consumerLimit, nil
}
