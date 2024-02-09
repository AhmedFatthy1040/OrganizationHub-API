package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
)

func BearerTokenAuth(userRepository *repository.UserRepository, organizationRepository *repository.OrganizationRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip authentication for signup, signin, and refresh-token routes
        if c.Request.URL.Path == "/users/signup" || c.Request.URL.Path == "/users/signin" || c.Request.URL.Path == "/users/refresh-token" {
            c.Next()
            return
        }

        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            return
        }

        token := tokenParts[1]

        // Verify token validity
        userID, err := verifyToken(token, userRepository)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Check if the user has access to the organization
        organizationID := c.Param("organization_id")
        if !organizationRepository.IsUserMemberOfOrganization(context.Background(), userID, organizationID) {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User does not have access to the organization"})
            return
        }

        // Set user ID in context for further use
        c.Set("user_id", userID)

        c.Next()
    }
}

func verifyToken(token string, userRepository *repository.UserRepository) (string, error) {
    // Retrieve the user by token from the repository
    user, err := userRepository.GetUserByToken(context.Background(), token)
    if err != nil {
        return "", err
    }

    // Return the user ID if found
    return user.ID.Hex(), nil
}
