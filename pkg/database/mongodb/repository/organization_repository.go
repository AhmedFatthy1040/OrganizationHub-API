package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/models"
)

type OrganizationRepository struct {
    collection *mongo.Collection
}

func NewOrganizationRepository(database *mongo.Database) *OrganizationRepository {
    return &OrganizationRepository{
        collection: database.Collection("organizations"),
    }
}

func (or *OrganizationRepository) CreateOrganization(ctx context.Context, org *models.Organization) error {
    org.CreatedAt = time.Now()
    org.UpdatedAt = time.Now()
    
    _, err := or.collection.InsertOne(ctx, org)
    if err != nil {
        log.Println("Error inserting organization:", err)
        return err
    }
    
    return nil
}

func (or *OrganizationRepository) GetOrganizationByID(ctx context.Context, id primitive.ObjectID) (*models.Organization, error) {
    var org models.Organization
    err := or.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&org)
    if err != nil {
        log.Println("Error getting organization by ID:", err)
        return nil, err
    }
    return &org, nil
}

func (or *OrganizationRepository) UpdateOrganization(ctx context.Context, id primitive.ObjectID, updatedOrg *models.Organization) error {
    updatedOrg.UpdatedAt = time.Now()
    _, err := or.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedOrg})
    if err != nil {
        log.Println("Error updating organization:", err)
        return err
    }
    return nil
}

func (or *OrganizationRepository) DeleteOrganization(ctx context.Context, id primitive.ObjectID) error {
    _, err := or.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        log.Println("Error deleting organization:", err)
        return err
    }
    return nil
}