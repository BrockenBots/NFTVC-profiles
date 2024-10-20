package middlewares

import (
	"net/http"
	"nftvc-profile/pkg/client"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type MiddlewareManager struct {
	log        logger.Logger
	cfg        *config.Config
	authClient *client.AuthClient
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config, authClient *client.AuthClient) *MiddlewareManager {
	return &MiddlewareManager{log: log, cfg: cfg, authClient: authClient}
}

func (m *MiddlewareManager) CORS() echo.MiddlewareFunc {
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}

	return middleware.CORSWithConfig(corsConfig)
}

func (m *MiddlewareManager) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		if accessToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
		}

		res, err := m.authClient.VerifyToken(accessToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: %"})
		}

		c.Set("accountId", res.AccountId)
		c.Set("token", accessToken)
		return next(c)
	}
}
