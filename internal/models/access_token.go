package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessToken struct {
	jwt.RegisteredClaims
	GUID string `json:"guid"`
}

func (t *AccessToken) Encode(secretKey string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, *t)
	tokenString, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *AccessToken) IntID() (int64, error) {
	n, err := strconv.ParseInt(t.ID, 10, 64)
	return n, err
}

func NewAccessToken(guid string, jti int64, expiresAt time.Time) *AccessToken {
	return &AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: fmt.Sprint(jti),
			ExpiresAt: jwt.NewNumericDate(
				expiresAt,
			),
		},
		GUID: guid,
	}
}

func NewAccessTokenFromString(token, secretKey string) (*AccessToken, error) {
	var claims AccessToken
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot parse jwt token %w", err)
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	return &claims, nil
}
