package service

import "math/rand"

type TransactionService interface {
	GenerateTransactionID() string
}

type transactionService struct{}

func NewTransactionService() TransactionService {
	return &transactionService{}
}

func (s *transactionService) GenerateTransactionID() string {
	prefix := "TRX"

	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randomPart := make([]byte, 9)
	for i := 0; i < 9; i++ {
		randomPart[i] = characters[rand.Intn(len(characters))]
	}

	return prefix + string(randomPart)
}
