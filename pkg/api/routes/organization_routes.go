package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/handlers"
)

// SetupOrganizationRoutes defines organization-related routes.
func SetupOrganizationRoutes(router *gin.Engine, organizationHandler *handlers.OrganizationHandler) {
	// Define a group for organization routes
	organizationRoutes := router.Group("/organizations")

	// Define routes for creating, reading, updating, and deleting organizations
	organizationRoutes.POST("/", organizationHandler.CreateOrganization)
	organizationRoutes.GET("/:id", organizationHandler.GetOrganizationByID)
	organizationRoutes.PUT("/:id", organizationHandler.UpdateOrganization)
	organizationRoutes.DELETE("/:id", organizationHandler.DeleteOrganization)

	// Define route for getting all organizations
	organizationRoutes.GET("/", organizationHandler.GetAllOrganizations)

	// Define route for inviting users to organizations
	organizationRoutes.POST("/:id/invite", organizationHandler.InviteUserToOrganization)
}
