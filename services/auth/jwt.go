package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Waris-Shaik/todo/types"
	"github.com/Waris-Shaik/todo/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(userID int) (string, error) {
	godotenv.Load()
	secret := os.Getenv("JWT_SECRET_KEY")
	if len(secret) == 0 {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}
	secretBytes := []byte(secret)

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiredAt": time.Now().Add(15 * time.Minute).Unix(), // Use Unix timestamp
	})

	// Sign token
	tokenString, err := token.SignedString(secretBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookies
		tokenString := getTokenFromCookie(r)
		if tokenString == "" {
			log.Println("Unauthorized: Empty token")
			permissionDenied(w)
			return
		}

		// Validate token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("Failed to validate JWT token:", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("Invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userIDFloat, ok := claims["userID"].(float64)
		if !ok {
			log.Println("Failed to assert userID type")
			permissionDenied(w)
			return
		}
		userID := int(userIDFloat)

		// Check token expiration
		expiredAtFloat, ok := claims["expiredAt"].(float64)
		if !ok {
			log.Println("Failed to assert expiredAt type")
			permissionDenied(w)
			return
		}
		if time.Unix(int64(expiredAtFloat), 0).Before(time.Now()) {
			log.Println("Token has expired")
			permissionDenied(w)
			return
		}

		user, err := store.GetUserByID(userID)
		if err != nil {
			log.Println("Failed to get user by ID:", err)
			permissionDenied(w)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserKey, user.ID)
		handlerFunc(w, r.WithContext(ctx))
	}
}

func getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("Error in getTokenFromCookie:", err)
		return ""
	}
	return cookie.Value
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("please login"))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return secret key
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		log.Println("Failed to get userID from context or userID is not int")
		return -1
	}
	return userID
}
