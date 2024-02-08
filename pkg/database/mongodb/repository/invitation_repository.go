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

type InvitationRepository struct {
    collection *mongo.Collection
}

func NewInvitationRepository(database *mongo.Database) *InvitationRepository {
    return &InvitationRepository{
        collection: database.Collection("invitations"),
    }
}

func (ir *InvitationRepository) CreateInvitation(ctx context.Context, invitation *models.Invitation) error {
    invitation.CreatedAt = time.Now()
    
    _, err := ir.collection.InsertOne(ctx, invitation)
    if err != nil {
        log.Println("Error inserting invitation:", err)
        return err
    }
    
    return nil
}

func (ir *InvitationRepository) GetInvitationByID(ctx context.Context, id primitive.ObjectID) (*models.Invitation, error) {
    var invitation models.Invitation
    err := ir.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&invitation)
    if err != nil {
        log.Println("Error getting invitation by ID:", err)
        return nil, err
    }
    return &invitation, nil
}

func (ir *InvitationRepository) UpdateInvitation(ctx context.Context, id primitive.ObjectID, updatedInvitation *models.Invitation) error {
    _, err := ir.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedInvitation})
    if err != nil {
        log.Println("Error updating invitation:", err)
        return err
    }
    return nil
}

func (ir *InvitationRepository) DeleteInvitation(ctx context.Context, id primitive.ObjectID) error {
    _, err := ir.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        log.Println("Error deleting invitation:", err)
        return err
    }
    return nil
}