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

type ConsumerHandler struct {
	consumerUsecase usecase.ConsumerUsecase
}

func NewConsumerHandler(consumerUsecase usecase.ConsumerUsecase) *ConsumerHandler {
	return &ConsumerHandler{
		consumerUsecase: consumerUsecase,
	}
}

func (h *ConsumerHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.NewError(model.ErrBadRequest, err)
		apiErr := httperror.FromError(appErr)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message, Details: err.Error()})
		return
	}

	auth, err := h.consumerUsecase.Login(c.Request.Context(), req.PhoneNumber, req.Password)
	if err != nil {
		apiErr := httperror.FromError(err)
		var details interface{}
		var modelErr model.Error
		if errors.As(err, &modelErr) {
			details = modelErr.AppError().Error()
		}
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message, Details: details})
		return
	}

	c.JSON(http.StatusOK, auth)
}

func (h *ConsumerHandler) GetByNIK(c *gin.Context) {
	nik := c.Param("nik")
	if nik == "" {
		appErr := model.NewError(model.ErrBadRequest, errors.New("NIK is required"))
		apiErr := httperror.FromError(appErr)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message})
		return
	}

	phoneNumber, err := util.GetUserPhoneNumberFromContext(c)
	if err != nil {
		appErr := model.NewError(model.ErrUnauthorized, err)
		apiErr := httperror.FromError(appErr)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message})
		return
	}

	consumer, err := h.consumerUsecase.GetByNIK(c.Request.Context(), phoneNumber, nik)
	if err != nil {
		apiErr := httperror.FromError(err)
		var details interface{}
		var modelErr model.Error
		if errors.As(err, &modelErr) {
			details = modelErr.AppError().Error()
		}
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message, Details: details})
		return
	}

	c.JSON(http.StatusOK, consumer)
}
