package repository

import (
	"context"
	"database/sql"

	"github.com/glennprays/xyz-fin/internal/app/model"
)

type ConsumerLimitRepository interface {
	FindByNIK(ctx context.Context, nik string) ([]*model.ConsumerLimit, error)
	FindByNIKAndTenor(ctx context.Context, tx *sql.Tx, nik string, tenor int) (*model.ConsumerLimit, error)
}

type consumerLimitRepository struct {
	db *sql.DB
}

func NewConsumerLimitRepository(db *sql.DB) ConsumerLimitRepository {
	return &consumerLimitRepository{db: db}
}

func (r *consumerLimitRepository) FindByNIK(ctx context.Context, nik string) ([]*model.ConsumerLimit, error) {
	var consumerLimits []*model.ConsumerLimit
	query := `SELECT consumer_nik, tenor, limit_amount, created_at, updated_at
			  FROM consumer_limits WHERE nik = $1 ORDER BY tenor ASC`

	rows, err := r.db.QueryContext(ctx, query, nik)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		consumerLimit := &model.ConsumerLimit{}
		if err := rows.Scan(&consumerLimit.ConsumerNIK, &consumerLimit.Tenor, &consumerLimit.LimitAmount, &consumerLimit.CreatedAt, &consumerLimit.UpdatedAt); err != nil {
			return nil, err
		}
		consumerLimits = append(consumerLimits, consumerLimit)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return consumerLimits, nil
}

func (r *consumerLimitRepository) FindByNIKAndTenor(ctx context.Context, tx *sql.Tx, nik string, tenor int) (*model.ConsumerLimit, error) {
	consumerLimit := &model.ConsumerLimit{}
	query := `SELECT consumer_nik, tenor, limit_amount, created_at, updated_at
        FROM consumer_limits WHERE consumer_nik = $1 AND tenor = $2`

	err := tx.QueryRowContext(ctx, query, nik, tenor).Scan(&consumerLimit.ConsumerNIK, &consumerLimit.Tenor, &consumerLimit.LimitAmount, &consumerLimit.CreatedAt, &consumerLimit.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return consumerLimit, nil
}
