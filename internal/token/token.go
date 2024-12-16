package token

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Newella-HQ/newella-backend/internal/model"
)

func ParseAccessToken(signedToken, signingKey string) (*model.NewellaJWTToken, error) {
	token, err := jwt.ParseWithClaims(signedToken, &model.NewellaJWTToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ivalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("can't parse jwt token: %w", err)
	}

	accessToken, ok := token.Claims.(*model.NewellaJWTToken)
	if !ok {
		return nil, fmt.Errorf("can't cast to newella jwt token: %w", err)
	}

	return accessToken, nil
}
