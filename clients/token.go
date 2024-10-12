package clients

import (
	"cms/models"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type tokenMaker struct {
	log    Logger
	config Config
}

func NewTokenMaker(log Logger, config Config) TokenMaker {
	return &tokenMaker{
		log:    log,
		config: config,
	}
}

type TokenMaker interface {
	GenerateTokenPair(sessionID uint) (models.Token, error)
	ValidateToken(tokenString string) (*models.Claims, error)
}

func (t *tokenMaker) GenerateTokenPair(sessionId uint) (models.Token, error) {
	var token models.Token
	var err error
	expirationTime := time.Now().Add(time.Duration(t.config.GetAccessTokenValidity()) * time.Hour)
	claims := &models.Claims{
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token.AccessToken, err = t.generateToken(claims)
	if err != nil {
		return models.Token{}, err
	}

	expirationTime = time.Now().Add(time.Duration(t.config.GetRefreshTokenValidity()) * time.Hour)
	claims = &models.Claims{
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token.RefreshToken, err = t.generateToken(claims)
	if err != nil {
		return models.Token{}, err
	}

	return token, nil
}

func (t *tokenMaker) generateToken(claims *models.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	JwtKey := []byte(t.config.GetJWTSecret())
	tokenString, err := token.SignedString(JwtKey)
	return tokenString, err
}

func (t *tokenMaker) ValidateToken(tokenString string) (*models.Claims, error) {
	if tokenString == "" {
		return nil, status.Error(codes.NotFound, "Token is empty")
	}

	tokens := strings.Split(tokenString, " ")
	if len(tokens) != 2 {
		return nil, status.Error(codes.NotFound, "Token is empty")
	}

	if tokens[0] != "Bearer" {
		return nil, status.Error(codes.NotFound, "Token is empty")
	}

	token := tokens[1]
	claims := &models.Claims{}
	parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.config.GetJWTSecret()), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, status.Error(codes.Unauthenticated, "Invalid token signature")
		}
		return nil, status.Error(codes.Unauthenticated, "Token is expired")
	}
	if !parseToken.Valid {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}

	return claims, nil
}
