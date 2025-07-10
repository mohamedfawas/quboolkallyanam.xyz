package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	customerrors "github.com/mohamedfawas/quboolkallyanam.xyz/pkg/errors"
)

type JWTConfig struct {
	SecretKey          string
	AccessTokenMinutes int
	RefreshTokenDays   int
	Issuer             string
}

type AppClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	config JWTConfig
}

func NewJWTManager(cfg JWTConfig) *JWTManager {
	return &JWTManager{config: cfg}
}

func (j *JWTManager) GenerateAccessToken(userID, role string) (string, error) {
	return j.generateToken(userID, role, time.Duration(j.config.AccessTokenMinutes)*time.Minute)
}

func (j *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	return j.generateToken(userID, "", time.Duration(j.config.RefreshTokenDays)*24*time.Hour)
}

func (j *JWTManager) generateToken(userID, role string, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := AppClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

func (j *JWTManager) VerifyToken(tokenStr string) (*AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customerrors.ErrInvalidToken
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, customerrors.ErrExpiredToken
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, customerrors.ErrTokenNotActive
		default:
			return nil, customerrors.ErrInvalidToken
		}
	}

	claims, ok := token.Claims.(*AppClaims)
	if !ok || !token.Valid {
		return nil, customerrors.ErrInvalidToken
	}

	if claims.Issuer != j.config.Issuer {
		return nil, customerrors.ErrInvalidToken
	}

	return claims, nil
}

func (j *JWTManager) ExtractUserIDAndRole(tokenStr string) (string, string, error) {
	claims, err := j.VerifyToken(tokenStr)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}
