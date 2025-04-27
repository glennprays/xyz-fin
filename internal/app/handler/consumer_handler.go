package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glennprays/xyz-fin/internal/app/httperror"
	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/glennprays/xyz-fin/internal/app/usecase"
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
