package routes

import (
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/handlers"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine, userRepository *repository.UserRepository) {
    // Initialize user handler
    userHandler := handlers.NewUserHandler(userRepository)

    // Define user-related routes
    userRoutes := router.Group("/users")
    {
        userRoutes.POST("/signup", userHandler.Signup)
        userRoutes.POST("/signin", userHandler.Signin)
        userRoutes.POST("/refresh-token", userHandler.RefreshToken)
    }
}
