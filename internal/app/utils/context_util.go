package util

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/glennprays/xyz-fin/internal/app/middleware"
	"github.com/glennprays/xyz-fin/internal/app/model"
)

func GetUserPhoneNumberFromContext(c *gin.Context) (string, error) {
	phoneNumber, exists := c.Get(middleware.ContextUserPhoneNumber)
	if !exists {
		return "", model.NewError(model.ErrUnauthorized, errors.New("user phone number not found in context (unauthorized)"))
	}

	phoneNumberStr, ok := phoneNumber.(string)
	if !ok {
		return "", model.NewError(model.ErrUnauthorized, errors.New("user phone number is not a string"))
	}

	return phoneNumberStr, nil
}
