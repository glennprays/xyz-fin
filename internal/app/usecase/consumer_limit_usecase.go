package usecase

import (
	"context"
	"errors"

	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/glennprays/xyz-fin/internal/app/repository"
)

type ConsumerLimitUsecase interface {
	GetLimitsByNIK(ctx context.Context, phoneNumber string, nik string) ([]model.ConsumerLimitResponse, error)
}

type consumerLimitUsecase struct {
	consumerRepo      repository.ConsumerRepository
	consumerLimitRepo repository.ConsumerLimitRepository
}

func NewConsumerLimitUsecase(
	consumerRepo repository.ConsumerRepository,
	consumerLimitRepo repository.ConsumerLimitRepository,
) ConsumerLimitUsecase {
	return &consumerLimitUsecase{
		consumerRepo:      consumerRepo,
		consumerLimitRepo: consumerLimitRepo,
	}
}

func (u *consumerLimitUsecase) GetLimitsByNIK(ctx context.Context, phoneNumber string, nik string) ([]model.ConsumerLimitResponse, error) {
	consumer, err := u.consumerRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		appErr := errors.New("consumer not found")
		return nil, model.NewError(model.ErrNotFound, appErr)
	}

	if consumer == nil {
		appErr := errors.New("consumer not found")
		return nil, model.NewError(model.ErrNotFound, appErr)
	}

	if consumer.NIK != nik {
		appErr := errors.New("nik does not match")
		return nil, model.NewError(model.ErrUnauthorized, appErr)
	}

	limits, err := u.consumerLimitRepo.FindByNIK(ctx, nik)
	if err != nil {
		appErr := errors.New("failed to find consumer limits")
		return nil, model.NewError(model.ErrInternalFailure, appErr)
	}

	var limitResponses []model.ConsumerLimitResponse
	for _, limit := range limits {
		limitResponses = append(limitResponses, model.ConsumerLimitResponse{
			LimitAmount: limit.LimitAmount,
			Tenor:       limit.Tenor,
		})
	}

	return limitResponses, nil
}
