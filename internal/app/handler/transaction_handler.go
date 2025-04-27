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

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req model.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := model.NewError(model.ErrBadRequest, err)
		apiErr := httperror.FromError(appErr)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message, Details: err.Error()})
		return
	}

	nik, err := util.GetUserPhoneNumberFromContext(c)
	if err != nil {
		apiErr := httperror.FromError(err)
		c.JSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message, Details: err.Error()})
		return
	}

	transaction, err := h.transactionUsecase.CreateTransaction(c.Request.Context(), nik, &req)
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

	c.JSON(http.StatusOK, transaction)
}
