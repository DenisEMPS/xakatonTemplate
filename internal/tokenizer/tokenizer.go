package tokenizer

import (
	"errors"
	"time"
	"xaxaton/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenInvalid = errors.New("access token invalid")
	ErrTokenExpired = errors.New("access token expired")
)

type TokenStruct struct {
	AccessToken string `json:"access_token"`
}

type Tokenizer interface {
	GenerateAccessTokenJWT(userID uuid.UUID) (*string, error)
	VerifyAccessTokenJWT(tokenString string) (jwt.MapClaims, error)
}

type tokenizer struct {
	tokenIssuer       string
	jwtSecret         []byte
	accessTokenExpire time.Duration
}

func New(cfg config.Tokenizer) Tokenizer {
	return &tokenizer{
		tokenIssuer:       cfg.IssuerName,
		jwtSecret:         []byte(cfg.AccessTokenSecretKey),
		accessTokenExpire: cfg.AccessTokenExpire,
	}
}

func (t *tokenizer) GenerateAccessTokenJWT(userID uuid.UUID) (*string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":     t.tokenIssuer,
		"exp":     time.Now().Add(t.accessTokenExpire).Unix(),
		"iat":     time.Now().Unix(),
		"user_id": userID,
	})

	accessToken, err := claims.SignedString(t.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (t *tokenizer) VerifyAccessTokenJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrTokenInvalid
		}
		return t.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
