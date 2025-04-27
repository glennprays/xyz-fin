package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glennprays/xyz-fin/internal/app/httperror"
	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/glennprays/xyz-fin/internal/app/usecase"
	util "github.com/glennprays/xyz-fin/internal/app/utils"
)

type ConsumerLimitHandler struct {
	consumerLimitUsecase usecase.ConsumerLimitUsecase
}

func NewConsumerLimitHandler(consumerLimitUsecase usecase.ConsumerLimitUsecase) *ConsumerLimitHandler {
	return &ConsumerLimitHandler{
		consumerLimitUsecase: consumerLimitUsecase,
	}
}

func (h *ConsumerLimitHandler) GetLimitsByNIK(c *gin.Context) {
	phoneNumber, err := util.GetUserPhoneNumberFromContext(c)
	if err != nil {
		apiErr := httperror.FromError(err)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message, Details: err.Error()})
		return
	}

	nik := c.Param("nik")
	if nik == "" {
		appErr := model.NewError(model.ErrBadRequest, errors.New("NIK is required"))
		apiErr := httperror.FromError(appErr)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message})
		return
	}

	limits, err := h.consumerLimitUsecase.GetLimitsByNIK(c.Request.Context(), phoneNumber, nik)
	if err != nil {
		apiErr := httperror.FromError(err)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message})
		return
	}

	c.JSON(http.StatusOK, limits)
}
