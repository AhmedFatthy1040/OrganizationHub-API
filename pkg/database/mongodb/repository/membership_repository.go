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

type MembershipRepository struct {
    collection *mongo.Collection
}

func NewMembershipRepository(database *mongo.Database) *MembershipRepository {
    return &MembershipRepository{
        collection: database.Collection("memberships"),
    }
}

func (mr *MembershipRepository) CreateMembership(ctx context.Context, membership *models.Membership) error {
    membership.CreatedAt = time.Now()
    membership.UpdatedAt = time.Now()
    
    _, err := mr.collection.InsertOne(ctx, membership)
    if err != nil {
        log.Println("Error inserting membership:", err)
        return err
    }
    
    return nil
}

func (mr *MembershipRepository) GetMembershipByID(ctx context.Context, id primitive.ObjectID) (*models.Membership, error) {
    var membership models.Membership
    err := mr.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&membership)
    if err != nil {
        log.Println("Error getting membership by ID:", err)
        return nil, err
    }
    return &membership, nil
}

func (mr *MembershipRepository) UpdateMembership(ctx context.Context, id primitive.ObjectID, updatedMembership *models.Membership) error {
    updatedMembership.UpdatedAt = time.Now()
    _, err := mr.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedMembership})
    if err != nil {
        log.Println("Error updating membership:", err)
        return err
    }
    return nil
}

func (mr *MembershipRepository) DeleteMembership(ctx context.Context, id primitive.ObjectID) error {
    _, err := mr.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        log.Println("Error deleting membership:", err)
        return err
    }
    return nil
}
