package auth

import (
	"errors"
	"template/app"
	"template/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}

const (
	_ACCESS_TOKEN_DURATION  = 5 * time.Minute
	_REFRESH_TOKEN_DURATION = 7 * 24 * time.Hour
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrUserIsNil          = errors.New("user is nil")
)

func ParseToken(authHeader string) (*jwt.Token, *Claims, error) {
	// Parse token
	if len(authHeader) <= len("Bearer ") || authHeader[:len("Bearer ")] != "Bearer " {
		zap.S().Debugf("token: %s", authHeader)
		return nil, nil, ErrInvalidTokenFormat
	}
	tokenString := authHeader[len("Bearer "):]
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(app.AccessKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	return token, &claims, nil
}

// GenerateTokens return a jwt access token and refresh token or an error
func GenerateTokens(user *model.User) (string, string, error) {
	if user == nil {
		return "", "", ErrUserIsNil
	}

	accessTokenClaims := &Claims{
		Username: user.Username,
		//		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_ACCESS_TOKEN_DURATION)),
			ID:        user.ID,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(app.AccessKey))
	if err != nil {
		zap.S().Debugf("Failed to generate access token err = %+v", err)
		return "", "", err
	}

	refreshTokenClaims := &Claims{
		Username: user.Username,
		//Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_REFRESH_TOKEN_DURATION)),
			ID:        user.ID,
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(app.RefreshKey))
	if err != nil {
		zap.S().Debugf("Failed to generate refresh token err = %+v", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
