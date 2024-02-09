package repository

import (
	"context"
	"errors"
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

func (or *OrganizationRepository) CreateOrganization(ctx context.Context, org *models.Organization) (primitive.ObjectID, error) {
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	result, err := or.collection.InsertOne(ctx, org)
	if err != nil {
		log.Println("Error inserting organization:", err)
		return primitive.NilObjectID, err
	}

	// Get the ID of the inserted organization
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Failed to get inserted ID")
		return primitive.NilObjectID, err
	}

	return insertedID, nil
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

func (or *OrganizationRepository) UpdateOrganization(ctx context.Context, id primitive.ObjectID, updatedOrg *models.Organization) (primitive.ObjectID, error) {
	updatedOrg.UpdatedAt = time.Now()

	result, err := or.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedOrg})
	if err != nil {
		log.Println("Error updating organization:", err)
		return primitive.NilObjectID, err
	}

	// Check if the document was found and updated
	if result.ModifiedCount == 0 {
		return primitive.NilObjectID, mongo.ErrNoDocuments
	}

	return id, nil
}

func (or *OrganizationRepository) DeleteOrganization(ctx context.Context, id primitive.ObjectID) error {
    _, err := or.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        log.Println("Error deleting organization:", err)
        return err
    }
    return nil
}

func (or *OrganizationRepository) AddMember(ctx context.Context, orgID primitive.ObjectID, member models.OrganizationMember) error {
	filter := bson.M{"_id": orgID}
	update := bson.M{"$push": bson.M{"organization_members": member}}

	_, err := or.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error adding member to organization:", err)
		return err
	}

	return nil
}

func (or *OrganizationRepository) GetAllOrganizations(ctx context.Context) ([]*models.Organization, error) {
    // Define a slice to store organizations
    var organizations []*models.Organization

    // Execute the find operation to retrieve all organizations
    cursor, err := or.collection.Find(ctx, bson.M{})
    if err != nil {
        log.Println("Error retrieving organizations:", err)
        return nil, err
    }

    // Iterate through the cursor and decode each document into an organization object
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var org models.Organization
        if err := cursor.Decode(&org); err != nil {
            log.Println("Error decoding organization:", err)
            continue
        }
        organizations = append(organizations, &org)
    }

    if err := cursor.Err(); err != nil {
        log.Println("Error iterating over organizations cursor:", err)
        return nil, err
    }

    return organizations, nil
}

func (or *OrganizationRepository) IsUserMemberOfOrganization(ctx context.Context, userID string, organizationID string) bool {
    // Convert organizationID to ObjectID
    oid, err := primitive.ObjectIDFromHex(organizationID)
    if err != nil {
        // Handle error
        return false
    }

    // Query the database to check if the user is a member of the organization
    filter := bson.M{
        "_id":              oid,
    }

    // Execute the query
    var result models.Organization
    err = or.collection.FindOne(ctx, filter).Decode(&result)
    if err != nil {
        
        return false
    }

    // Check if the user is a member of the organization
    for _, member := range result.Members {
        user, err := or.GetUserByEmail(ctx, member.Email)
        if err != nil {
            return false
        }

        if user.ID.String() == userID {
            return true
        }
    }

    return false
}

func (or *OrganizationRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    pipeline := bson.A{
        bson.M{"$match": bson.M{"email": email}},
    }

    // Perform aggregation
    cursor, err := or.collection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    // Decode the result
    var user models.User
    if cursor.Next(ctx) {
        if err := cursor.Decode(&user); err != nil {
            return nil, err
        }
        return &user, nil
    }

    // If no user found with the given email
    return nil, errors.New("user not found")
}