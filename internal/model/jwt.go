package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"access_token" db:"access_token"`   // signed NewellaJWTToken - in system | db - google AT
	RefreshToken string `json:"refresh_token" db:"refresh_token"` // from Google
}

type NewellaJWTToken struct {
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	Audience       string `json:"aud"`
	ExpirationTime int64  `json:"exp"`
	IssuedAt       int64  `json:"iat"`
	Issuer         string `json:"iss"`
}

func (n NewellaJWTToken) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(n.ExpirationTime, 0).UTC()), nil
}

func (n NewellaJWTToken) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(n.IssuedAt, 0).UTC()), nil
}

func (n NewellaJWTToken) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now().UTC()), nil
}

func (n NewellaJWTToken) GetIssuer() (string, error) {
	return n.Issuer, nil
}

func (n NewellaJWTToken) GetSubject() (string, error) {
	return n.UserID, nil
}

func (n NewellaJWTToken) GetAudience() (jwt.ClaimStrings, error) {
	return []string{n.Audience}, nil
}

type OAuthJWTToken struct {
	Audience        string `json:"aud"`
	ExpirationTime  int64  `json:"exp"`
	IssuedAt        int64  `json:"iat"`
	Issuer          string `json:"iss"`
	Subject         string `json:"sub"`
	AccessTokenHash string `json:"at_hash"`
	AuthorizedParty string `json:"azp"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	FamilyName      string `json:"family_name"`
	GivenName       string `json:"given_name"`
	Name            string `json:"name"`
	Picture         string `json:"picture"`
}

func (O OAuthJWTToken) Validate() (OAuthJWTToken, error) {
	if O.Audience == "" {
		return OAuthJWTToken{}, fmt.Errorf("empty aud")
	}
	if O.ExpirationTime == 0 {
		return OAuthJWTToken{}, fmt.Errorf("empty exp")
	}
	if O.IssuedAt == 0 {
		return OAuthJWTToken{}, fmt.Errorf("empty iat")
	}
	if O.Subject == "" {
		return OAuthJWTToken{}, fmt.Errorf("empty subject")
	}
	if O.AccessTokenHash == "" {
		return OAuthJWTToken{}, fmt.Errorf("empty at_hash")
	}
	if O.AuthorizedParty == "" {
		return OAuthJWTToken{}, fmt.Errorf("empty azp")
	}
	if O.Email == "" {
		return OAuthJWTToken{}, fmt.Errorf("empty email")
	}
	if O.Name == "" {
		O.Name = O.GivenName + " " + O.FamilyName
	}
	return O, nil
}

func (O OAuthJWTToken) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(O.ExpirationTime, 0).UTC()), nil
}

func (O OAuthJWTToken) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(O.IssuedAt, 0).UTC()), nil
}

func (O OAuthJWTToken) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now().UTC()), nil
}

func (O OAuthJWTToken) GetIssuer() (string, error) {
	return O.Issuer, nil
}

func (O OAuthJWTToken) GetSubject() (string, error) {
	return O.Subject, nil
}

func (O OAuthJWTToken) GetAudience() (jwt.ClaimStrings, error) {
	return []string{O.Audience}, nil
}
