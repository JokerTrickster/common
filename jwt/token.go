package jwt

/*
	JWT 생성, 검증, 파싱 로직
*/

import (
	"context"
	"fmt"
	"time"

	_error "github.com/JokerTrickster/common/error"
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	CreateTime int64  `json:"createTime"`
	UserID     uint   `json:"userID"`
	Email      string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken generates access and refresh tokens
func GenerateToken(email string, userID uint) (string, int64, string, int64, error) {
	now := time.Now()
	accessToken, accessTknExpiredAt, err := GenerateAccessToken(email, now, userID)
	if err != nil {
		return "", 0, "", 0, err
	}
	refreshToken, refreshTknExpiredAt, err := GenerateRefreshToken(email, now, userID)
	if err != nil {
		return "", 0, "", 0, err
	}
	return accessToken, accessTknExpiredAt, refreshToken, refreshTknExpiredAt, nil
}

// GenerateAccessToken generates an access token
func GenerateAccessToken(email string, now time.Time, userID uint) (string, int64, error) {
	expiredAt := now.Add(time.Hour * AccessTokenExpiredTime).Unix()
	claims := &JwtCustomClaims{
		CreateTime: now.Unix(),
		UserID:     userID,
		Email:      email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(AccessTokenSecretKey)
	if err != nil {
		return "", 0, err
	}
	return accessToken, expiredAt, nil
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(email string, now time.Time, userID uint) (string, int64, error) {
	expiredAt := now.Add(time.Hour * RefreshTokenExpiredTime).Unix()
	claims := &JwtCustomClaims{
		CreateTime: now.Unix(),
		UserID:     userID,
		Email:      email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString(RefreshTokenSecretKey)
	if err != nil {
		return "", 0, _error.CreateError(context.TODO(), string(_error.ErrBadToken), _error.Trace(), fmt.Sprintf("failed to generate refresh token: %v", err), string(_error.ErrFromInternal))
	}
	return refreshToken, expiredAt, nil
}

// VerifyToken verifies the validity of a JWT
func VerifyToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	if err != nil || !token.Valid {
		return _error.CreateError(context.TODO(), string(_error.ErrBadToken), _error.Trace(), "invalid token", string(_error.ErrFromClient))
	}
	return nil
}

// ParseToken extracts claims from a JWT
func ParseToken(tokenString string) (uint, string, error) {
	token, _ := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessTokenSecretKey, nil
	})
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return 0, "", _error.CreateError(context.TODO(), string(_error.ErrBadToken), _error.Trace(), "failed to extract claims", string(_error.ErrFromClient))
	}
	return claims.UserID, claims.Email, nil
}
