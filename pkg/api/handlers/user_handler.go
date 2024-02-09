package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/models"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/repository"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
    userRepository *repository.UserRepository
}

func NewUserHandler(userRepository *repository.UserRepository) *UserHandler {
    return &UserHandler{
        userRepository: userRepository,
    }
}

func (uh *UserHandler) Signup(c *gin.Context) {
    var user models.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash the user's password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    user.Password = string(hashedPassword)

    // Set up user creation timestamp
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    // Generate access token
    accessToken, err := utils.GenerateAccessToken(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    // Generate refresh token
    refreshToken := utils.GenerateRefreshToken()

    // Update user model with access token and refresh token
    user.AccessToken = accessToken
    user.RefreshToken = refreshToken

    // Save the user in the database
    if err := uh.userRepository.CreateUser(context.Background(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Respond with access token and refresh token
    c.JSON(http.StatusOK, gin.H{
        "message":        "User signed up successfully",
    })
}

func (uh *UserHandler) Signin(c *gin.Context) {
    // Bind the request body to the User model
    var user models.User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Find the user by email in the database
    foundUser, err := uh.userRepository.GetUserByEmail(context.Background(), user.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Verify the password
    if !utils.VerifyPassword(user.Password, foundUser.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Generate access token
    accessToken, err := utils.GenerateAccessToken(foundUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    // Generate refresh token
    refreshToken := utils.GenerateRefreshToken()

    // Update user model with access token and refresh token
    foundUser.AccessToken = accessToken
    foundUser.RefreshToken = refreshToken

    // Save the updated user in the database
    if err := uh.userRepository.UpdateUser(context.Background(), foundUser.ID, foundUser); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    // Respond with access token and refresh token
    c.JSON(http.StatusOK, gin.H{
        "message":        "User signed in successfully",
        "access_token":   accessToken,
        "refresh_token":  refreshToken,
    })
}

func (uh *UserHandler) RefreshToken(c *gin.Context) {
    // Bind the request body to the RefreshTokenRequest model
    var req models.User
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Validate the refresh token by its value
    userID, err := uh.userRepository.ValidateTokensByValue(context.Background(), req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
        return
    }

    // Generate new access token
    user := models.User{ID: userID}
    accessToken, err := utils.GenerateAccessToken(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    // Generate new refresh token
    refreshToken := utils.GenerateRefreshToken()

    // Save the new tokens in the database
    if err := uh.userRepository.SaveTokens(context.Background(), userID, accessToken, refreshToken); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save tokens"})
        return
    }

    // Respond with new access token and refresh token
    c.JSON(http.StatusOK, gin.H{
        "message":        "Tokens refreshed successfully",
        "access_token":   accessToken,
        "refresh_token":  refreshToken,
    })
}
