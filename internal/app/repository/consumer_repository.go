package repository

import (
	"database/sql"

	"github.com/glennprays/xyz-fin/internal/app/model"
)

type ConsumerRepository interface {
	FindByPhoneNumber(phoneNumber string) (*model.Consumer, error)
	FindByNIK(nik string) (*model.Consumer, error)
}

type consumerRepository struct {
	db *sql.DB
}

func NewConsumerRepository(db *sql.DB) ConsumerRepository {
	return &consumerRepository{db: db}
}

func (r *consumerRepository) FindByPhoneNumber(phoneNumber string) (*model.Consumer, error) {
	consumer := &model.Consumer{}
	query := `SELECT nik, phone_number, password_hash, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp_path, foto_selfie_path, created_at, updated_at 
  FROM consumers WHERE phone_number = $1`
	err := r.db.QueryRow(query, phoneNumber).Scan(&consumer.NIK, &consumer.PhoneNumber, &consumer.PasswordHash, &consumer.FullName, &consumer.LegalName, &consumer.TempatLahir, &consumer.TanggalLahir, &consumer.Gaji, &consumer.FotoKTPPath, &consumer.FotoSelfiePath, &consumer.CreatedAt, &consumer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return consumer, nil
}

func (r *consumerRepository) FindByNIK(nik string) (*model.Consumer, error) {
	consumer := &model.Consumer{}
	query := `SELECT nik, phone_number, password_hash, full_name, legal_name, tempat_lahir, tanggal_lahir, gaji, foto_ktp_path, foto_selfie_path, created_at, updated_at 
  FROM consumers WHERE nik = $1`
	err := r.db.QueryRow(query, nik).Scan(&consumer.NIK, &consumer.PhoneNumber, &consumer.PasswordHash, &consumer.FullName, &consumer.LegalName, &consumer.TempatLahir, &consumer.TanggalLahir, &consumer.Gaji, &consumer.FotoKTPPath, &consumer.FotoSelfiePath, &consumer.CreatedAt, &consumer.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return consumer, nil
}
