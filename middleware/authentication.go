package middleware

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct { 
	secretKey []byte
	tokenTTL  time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

// NewAuth initializes a new instance of the Auth middleware
func NewAuth(secretKey string, tokenTTL time.Duration) (*Auth, error) {
	// Validate and store the secret key
	if len(secretKey) == 0 {
		return nil, errors.New("secret key is required")
	}
	secretKeyBytes := []byte(secretKey)

	return &Auth{
		secretKey: secretKeyBytes,
		tokenTTL:  tokenTTL,
	}, nil
}

// GenerateToken generates a new JWT token for the provided user ID
func (a *Auth) GenerateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.tokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(a.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyToken verifies the provided JWT token and returns the user ID if the token is valid
func (a *Auth) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.UserID, nil
}

// HashPassword hashes the provided password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// VerifyPassword verifies the provided password against the hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
