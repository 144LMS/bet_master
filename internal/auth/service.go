package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secret          string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		secret:          secret,
		accessDuration:  15 * time.Minute,
		refreshDuration: 7 * 24 * time.Hour,
	}
}

func (s *AuthService) GenerateTokens(id uint, role string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  id,
		"role": role,
		"exp":  time.Now().Add(s.accessDuration).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(s.refreshDuration).Unix(),
	})

	accessSigned, err := accessToken.SignedString([]byte(s.secret))
	if err != nil {
		return "", "", err
	}

	refreshSigned, err := refreshToken.SignedString([]byte(s.secret + "refresh"))
	if err != nil {
		return "", "", err
	}

	return accessSigned, refreshSigned, nil
}

func (s *AuthService) ValidateTokens(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return 0, "", jwt.ErrTokenExpired
	}

	userID, ok := claims["sub"].(uint)
	if !ok {
		return 0, "", jwt.ErrTokenInvalidId
	}

	role, _ := claims["role"].(string)

	return userID, role, nil
}
