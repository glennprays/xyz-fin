package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glennprays/xyz-fin/internal/app/httperror"
	"github.com/glennprays/xyz-fin/internal/app/model"
	"github.com/glennprays/xyz-fin/pkg/auth"
)

const (
	ContextUserIDKey        = "userID"
	ContextUserRoleKey      = "userRole"
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "Bearer"
)

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := m.extractBearerToken(c)
		if err != nil {
			apiErr := httperror.FromError(err)
			c.AbortWithStatusJSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message})
			return
		}

		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			var (
				message = "Invalid token"
				svcErr  = model.ErrUnauthorized
			)

			switch {
			case errors.Is(err, auth.ErrTokenExpired):
				message = "Token has expired"
			case errors.Is(err, auth.ErrTokenInvalid), errors.Is(err, auth.ErrMissingClaims):
				message = "Invalid token"
			default:
				message = "Internal error during token validation"
			}

			appErr := model.NewError(svcErr, errors.New(message))
			apiErr := httperror.FromError(appErr)
			c.AbortWithStatusJSON(apiErr.Status, model.ErrorResponse{Message: apiErr.Message})
			return
		}

		c.Set(ContextUserIDKey, claims.PhoneNumber)

		c.Next()
	}
}

func (m *AuthMiddleware) extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(AuthorizationHeaderKey)
	if authHeader == "" {
		appErr := errors.New("authorization header required")
		return "", model.NewError(model.ErrUnauthorized, appErr)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], AuthorizationTypeBearer) {
		appErr := errors.New("invalid authorization header format (must be Bearer token)")
		return "", model.NewError(model.ErrUnauthorized, appErr)
	}

	return parts[1], nil
}
