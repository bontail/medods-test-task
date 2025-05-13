package models

import (
	"fmt"
	"net/netip"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RefreshToken struct {
	Id          int64      `json:"id"`
	UserGUID    string     `json:"user_guid"`
	SecretValue string     `json:"secret_value"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	BlockedAt   time.Time  `json:"blocked_at"`
	UserAgent   string     `json:"user_agent"`
	IP          netip.Addr `json:"ip"`
}

func NewRefreshToken(guid, secretValue string, createdAt, ExpiresAt time.Time, userAgent string, ip netip.Addr) (RefreshToken, error) {
	hashSecretValue, err := bcrypt.GenerateFromPassword([]byte(secretValue), bcrypt.DefaultCost)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("cannot create refresh token: %w", err)
	}

	return RefreshToken{
		UserGUID:    guid,
		SecretValue: string(hashSecretValue),
		CreatedAt:   createdAt,
		ExpiresAt:   ExpiresAt,
		UserAgent:   userAgent,
		IP:          ip,
	}, nil
}

func (t *RefreshToken) CompareSecretValue(value string) bool {
	return bcrypt.CompareHashAndPassword([]byte(t.SecretValue), []byte(value)) == nil
}
