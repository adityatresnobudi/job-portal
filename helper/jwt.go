package helper

import (
	"time"

	"github.com/adityatresnobudi/job-portal/dto"
	"github.com/adityatresnobudi/job-portal/shared"
	"github.com/golang-jwt/jwt/v5"
)

var APPLICATION_NAME = "Library"
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("3 permen milkita sama dengan 1 gelas susu")

type JWTClaims struct {
	jwt.RegisteredClaims
	UserId uint `json:"id"`
}

func AuthorizedJWT(claims JWTClaims, user dto.UserPayload) (string, error) {
	claims.Issuer = APPLICATION_NAME
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	claims.UserId = user.ID

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

	generateToken, err := token.SignedString([]byte(JWT_SIGNATURE_KEY))
	if err != nil {
		return "", err
	}

	return generateToken, nil
}

func ValidateJWT(generateToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(generateToken, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, shared.ErrInvalidToken
		}

		return JWT_SIGNATURE_KEY, nil
	})
}
