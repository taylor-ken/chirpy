package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)

	claims := &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   string(userID.String()),
		Issuer:    "chirpy",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		var nilUUID uuid.UUID
		return nilUUID, err
	}
	claims := token.Claims.(*jwt.RegisteredClaims)
	subject, err := uuid.Parse(claims.Subject)
	if err != nil {
		var nilUUID uuid.UUID
		return nilUUID, err
	}
	return subject, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	// Check if the header starts with "Bearer "
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return strings.TrimSpace(authHeader[7:]), nil
	}

	return "", errors.New("invalid authorization header format")
}

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.New("error generating random bytes")
	}
	token := hex.EncodeToString(randomBytes)
	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not found")
	}

	// Check if the header starts with "Bearer "
	if len(authHeader) > 7 && authHeader[:7] == "ApiKey " {
		return strings.TrimSpace(authHeader[7:]), nil
	}

	return "", errors.New("invalid authorization header format")
}
