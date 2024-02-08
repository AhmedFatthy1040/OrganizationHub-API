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

    // Create the user in the database
    if err := uh.userRepository.CreateUser(context.Background(), &user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Generate access token
    accessToken, err := utils.GenerateAccessToken(&user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
        return
    }

    // Generate refresh token
    refreshToken := utils.GenerateRefreshToken()

    // Respond with access token and refresh token
    c.JSON(http.StatusOK, gin.H{
        "message":        "User signed up successfully",
        "access_token":   accessToken,
        "refresh_token":  refreshToken,
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

    // Respond with access token and refresh token
    c.JSON(http.StatusOK, gin.H{
        "message":        "User signed in successfully",
        "access_token":   accessToken,
        "refresh_token":  refreshToken,
    })
}