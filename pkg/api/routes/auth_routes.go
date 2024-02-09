package routes

import (
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/handlers"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/middleware"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/repository"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, userRepository *repository.UserRepository, organizationRepository *repository.OrganizationRepository) {
    // Initialize organization handler
    organizationHandler := handlers.NewOrganizationHandler(organizationRepository)

    // Create a new router group for authenticated routes
    authRoutes := router.Group("/auth")
    authRoutes.Use(middleware.BearerTokenAuth(userRepository, organizationRepository))

    {
        // Organization routes
        authRoutes.POST("/organization", organizationHandler.CreateOrganization)
        authRoutes.GET("/organization/:organization_id", organizationHandler.GetOrganizationByID)
        authRoutes.GET("/organization", organizationHandler.GetAllOrganizations)
        authRoutes.PUT("/organization/:organization_id", organizationHandler.UpdateOrganization)
        authRoutes.DELETE("/organization/:organization_id", organizationHandler.DeleteOrganization)
        authRoutes.POST("/organization/:organization_id/invite", organizationHandler.InviteUserToOrganization)
    }
}
