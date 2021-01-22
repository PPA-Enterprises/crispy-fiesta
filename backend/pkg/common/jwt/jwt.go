package jwt

import (
	"fmt"
	"strings"
	"time"

	"github.com/ribice/gorsk"
	"github.com/dgrijalva/jwt-go"
)

var minSecretLen = 128

func New(algo, secret string, ttlMinutes, minSecretLength int) (Service, error) {
	if minSecretLength > 0 {
		minSecretLen = minSecretLength
	}
	if len(secret) < minSecretLen {
		return Service{}, fmt.Errorf("jwt secret length is %v, which is less than required %v", len(secret), minSecretLen)
	}
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		return Service{}, fmt.Errorf("invalid jwt signing method: %s", algo)
	}

	return Service{
		key:  []byte(secret),
		algo: signingMethod,
		ttl:  time.Duration(ttlMinutes) * time.Minute,
	}, nil
}

type Service struct {
	key []byte
	ttl time.Duration
	algo jwt.SigningMethod
}

func (s Service) ParseToken(authHeader string) (*jwt.Token, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, gorsk.ErrGeneric
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if s.algo != token.Method {
			return nil, gorsk.ErrGeneric
		}
		return s.key, nil
	})

}

func (s Service) GenerateToken(u gorsk.User) (string, error) {
	return jwt.NewWithClaims(s.algo, jwt.MapClaims{
		"_id":  u.ID.hex(),
		"e":   u.Email,
		"r":   u.Role.AccessLevel,
		"exp": time.Now().Add(s.ttl).Unix(),
	}).SignedString(s.key)

}
