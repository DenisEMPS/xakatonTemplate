package middleware

import (
	"errors"
	"strings"
	"xaxaton/internal/tokenizer"

	"github.com/gofiber/fiber/v2"
)

const (
	AuthorizationHeader                = "Authorization"
	BearerScheme                       = "Bearer"
	UserID              userContextKey = "user_id"
)

type userContextKey string

type Middleware interface {
	JWT() fiber.Handler
	RecoverMiddleware() fiber.Handler
}

type middleware struct {
	tokenizer tokenizer.Tokenizer
}

func New(tok tokenizer.Tokenizer) Middleware {
	return &middleware{
		tokenizer: tok,
	}
}

func (m *middleware) RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				c.Context().Logger().Printf("%v", err)
				c.Context().SetStatusCode(fiber.ErrInternalServerError.Code)
			}
		}()

		return c.Next()
	}
}

func (m *middleware) JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken, err := extractToken(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		}

		claims, err := m.tokenizer.VerifyAccessTokenJWT(accessToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("unauthorized")
		}

		c.Locals(UserID, claims["user_id"])

		return c.Next()
	}
}

func extractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get(AuthorizationHeader)

	if authHeader == "" {
		return "", errors.New("auth header is empty")
	}

	authHeader = strings.TrimSpace(authHeader)
	tokenParts := strings.Split(authHeader, " ")

	if len(tokenParts) != 2 || !strings.EqualFold(tokenParts[0], BearerScheme) {
		return "", errors.New("invalid token format")
	}

	return tokenParts[1], nil
}
