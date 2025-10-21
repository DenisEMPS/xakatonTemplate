package middleware

import (
	"errors"
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
	tokenStruct := tokenizer.TokenStruct{}

	err := c.BodyParser(&tokenStruct)
	if err != nil {
		return "", err
	}

	if tokenStruct.AccessToken == "" {
		return "", errors.New("authorization header not provided")
	}

	return tokenStruct.AccessToken, nil
}
