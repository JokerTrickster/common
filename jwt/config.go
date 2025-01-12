package jwt

/*
	JWT 초기화 및 설정
*/
import (
	"github.com/labstack/echo/v4/middleware"
)

var (
	AccessTokenSecretKey  []byte
	RefreshTokenSecretKey []byte
)

var JwtConfig middleware.JWTConfig

const (
	AccessTokenExpiredTime  = 24     // hour
	RefreshTokenExpiredTime = 24 * 7 // hour
)

// InitJWT initializes the JWT secrets and configuration
func InitJWT() error {
	// Initialize JWT Secret Keys
	secret := "secret" // Replace this with a secure secret key
	AccessTokenSecretKey = []byte(secret)
	RefreshTokenSecretKey = []byte(secret)

	// Configure JWT Middleware
	JwtConfig = middleware.JWTConfig{
		SigningKey:  AccessTokenSecretKey,   // Signing key for token verification
		TokenLookup: "header:Authorization", // Where to look for the token
		AuthScheme:  "Bearer",               // Authorization header prefix
	}

	return nil
}
