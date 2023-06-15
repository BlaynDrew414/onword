package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Authorization represents the authorization middleware
type Authorization struct {
	allowedRoles []string
	secretKey    []byte
}

// NewAuthorization initializes a new instance of the Authorization middleware
func NewAuthorization(allowedRoles []string, secretKey string) (*Authorization, error) {
	if len(secretKey) == 0 {
		return nil, errors.New("secret key is required")
	}

	secretKeyBytes := []byte(secretKey)

	return &Authorization{
		allowedRoles: allowedRoles,
		secretKey:    secretKeyBytes,
	}, nil
}

// Middleware handles the authorization logic
func (a *Authorization) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify and parse the JWT token 
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return a.secretKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user's role is allowed
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !a.isRoleAllowed(claims["role"].(string)) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Continue to the next handler if authorized
		next.ServeHTTP(w, r)
	})
}

// isRoleAllowed checks if the user's role is allowed
func (a *Authorization) isRoleAllowed(role string) bool {
	for _, allowedRole := range a.allowedRoles {
		if role == allowedRole {
			return true
		}
	}
	return false
}
