package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type (
	Claims struct {
		SessionId uint
		jwt.RegisteredClaims
	}
)

type Token struct {
	AccessToken  string
	RefreshToken string
}
