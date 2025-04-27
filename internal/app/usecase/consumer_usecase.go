package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/glennprays/xyz-fin/internal/app/repository"
	"github.com/glennprays/xyz-fin/pkg/auth"
	"github.com/glennprays/xyz-fin/pkg/hasher"
)

type ConsumerUsecase interface {
	Login(ctx context.Context, phoneNumber string, password string) (*model.AuthResponse, error)
}

type consumerUsecase struct {
	repo       repository.ConsumerRepository
	jwtManager *auth.JWTManager
	hasher     hasher.Argon2idHasher
}

func NewConsumerUsecase(repo repository.ConsumerRepository, jwtManager *auth.JWTManager, hasher hasher.Argon2idHasher) ConsumerUsecase {
	return &consumerUsecase{
		repo:       repo,
		jwtManager: jwtManager,
		hasher:     hasher,
	}
}

func (u *consumerUsecase) Login(ctx context.Context, phoneNumber string, password string) (*model.AuthResponse, error) {
	consumer, err := u.repo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		appErr := errors.New("consumer not found")
		return nil, model.NewError(model.ErrNotFound, appErr)
	}

	if !u.hasher.Check(password, consumer.PasswordHash) {
		appErr := errors.New("invalid password")
		return nil, model.NewError(model.ErrUnauthorized, appErr)
	}

	accessToken, refreshToken, err := u.jwtManager.GenerateTokens(phoneNumber)
	if err != nil {
		appErr := fmt.Errorf("failed to generate authentication tokens: %w", err)
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
