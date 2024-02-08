// utils/tokens.go

package utils

import (
	"os"
	"time"

	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/models"
    
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// GenerateAccessToken generates an access token for the given user.
func GenerateAccessToken(user *models.User) (string, error) {
    // Define the expiration time for the token (1 hour)
    expirationTime := time.Now().Add(1 * time.Hour)

    // Create a new JWT token with the user's ID and expiration time
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "exp":     expirationTime.Unix(),
    })

    // Sign the token with a secret key
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// GenerateRefreshToken generates a refresh token for the given user.
func GenerateRefreshToken() string {
    // Generate a UUID as the refresh token
    refreshToken := uuid.New().String()
    return refreshToken
}

// VerifyPassword checks if the provided password matches the hashed password.
func VerifyPassword(plainPassword, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
    return err == nil
}