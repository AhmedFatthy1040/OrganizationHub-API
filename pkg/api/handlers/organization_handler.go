package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/models"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/repository"
)

type OrganizationHandler struct {
	organizationRepository *repository.OrganizationRepository
}

func NewOrganizationHandler(organizationRepository *repository.OrganizationRepository) *OrganizationHandler {
	return &OrganizationHandler{
		organizationRepository: organizationRepository,
	}
}

func (oh *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var organization models.Organization
	if err := c.BindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set up organization creation timestamp
	organization.CreatedAt = time.Now()
	organization.UpdatedAt = time.Now()

	// Create the organization in the database
	id, err := oh.organizationRepository.CreateOrganization(context.Background(), &organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	// Respond with the created organization ID
	c.JSON(http.StatusOK, gin.H{
		"organization_id": id.Hex(),
	})
}

func (oh *OrganizationHandler) GetOrganizationByID(c *gin.Context) {
	organizationID := c.Param("id")

	// Convert organization ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Retrieve the organization by ID from the database
	organization, err := oh.organizationRepository.GetOrganizationByID(context.Background(), objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	// Respond with the retrieved organization
	c.JSON(http.StatusOK, gin.H{
		"organization_id":   organization.ID.Hex(),
		"name":              organization.Name,
		"description":       organization.Description,
		"organization_members": organization.Members,
	})
}

func (oh *OrganizationHandler) GetAllOrganizations(c *gin.Context) {
	// Retrieve all organizations from the database
	organizations, err := oh.organizationRepository.GetAllOrganizations(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organizations"})
		return
	}

	// Respond with the retrieved organizations
	var response []gin.H
	for _, org := range organizations {
		response = append(response, gin.H{
			"organization_id":   org.ID.Hex(),
			"name":              org.Name,
			"description":       org.Description,
			"organization_members": org.Members,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (oh *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	organizationID := c.Param("id")

	// Convert organization ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var organization models.Organization
	if err := c.BindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set up organization update timestamp
	organization.UpdatedAt = time.Now()

	// Update the organization in the database
	id, err := oh.organizationRepository.UpdateOrganization(context.Background(), objectID, &organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	// Respond with the updated organization
	c.JSON(http.StatusOK, gin.H{
		"organization_id":   id.Hex(),
		"name":              organization.Name,
		"description":       organization.Description,
	})
}

func (oh *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	organizationID := c.Param("id")

	// Convert organization ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Delete the organization from the database
	if err := oh.organizationRepository.DeleteOrganization(context.Background(), objectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

func (oh *OrganizationHandler) InviteUserToOrganization(c *gin.Context) {
	// Parse organization ID from request parameters
	organizationID := c.Param("id")

	// Bind the request body to the InviteUserRequest model
	var inviteRequest models.Invitation
	if err := c.BindJSON(&inviteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate user email
	if inviteRequest.InvitedEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User email is required"})
		return
	}

	// Check if the organization exists
	orgID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}
	organization, err := oh.organizationRepository.GetOrganizationByID(c, orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	// // Check if the user making the request is an admin
	// if !userIsAdmin(c) {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can invite users"})
	// 	return
	// }

	// Invite the user to the organization (add user's email to organization members)
	organization.Members = append(organization.Members, models.OrganizationMember{
		Name:        "Test",
		Email:       inviteRequest.InvitedEmail,
		AccessLevel: "member",
	})

	// Update the organization in the database
	id, err := oh.organizationRepository.UpdateOrganization(c, orgID, organization)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite user to organization"})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User invited to organization successfully",
		"id": id.Hex(),
	})
}