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

type TokenRepository struct {
    collection *mongo.Collection
}

func NewTokenRepository(database *mongo.Database) *TokenRepository {
    return &TokenRepository{
        collection: database.Collection("tokens"),
    }
}

func (tr *TokenRepository) CreateToken(ctx context.Context, token *models.Token) error {
    token.CreatedAt = time.Now()
    
    _, err := tr.collection.InsertOne(ctx, token)
    if err != nil {
        log.Println("Error inserting token:", err)
        return err
    }
    
    return nil
}

func (tr *TokenRepository) GetTokenByID(ctx context.Context, id primitive.ObjectID) (*models.Token, error) {
    var token models.Token
    err := tr.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&token)
    if err != nil {
        log.Println("Error getting token by ID:", err)
        return nil, err
    }
    return &token, nil
}

func (tr *TokenRepository) UpdateToken(ctx context.Context, id primitive.ObjectID, updatedToken *models.Token) error {
    _, err := tr.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedToken})
    if err != nil {
        log.Println("Error updating token:", err)
        return err
    }
    return nil
}

func (tr *TokenRepository) DeleteToken(ctx context.Context, id primitive.ObjectID) error {
    _, err := tr.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        log.Println("Error deleting token:", err)
        return err
    }
    return nil
}