package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/models"
)

type UserRepository struct {
    collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
    return &UserRepository{
        collection: database.Collection("users"),
    }
}

func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    _, err := ur.collection.InsertOne(ctx, user)
    if err != nil {
        log.Println("Error inserting user:", err)
        return err
    }
    
    return nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
    var user models.User
    err := ur.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        log.Println("Error getting user by ID:", err)
        return nil, err
    }
    return &user, nil
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := ur.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("user with email %s not found", email)
        }
        log.Println("Error getting user by email:", err)
        return nil, err
    }
    return &user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, updatedUser *models.User) error {
    updatedUser.UpdatedAt = time.Now()
    _, err := ur.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedUser})
    if err != nil {
        log.Println("Error updating user:", err)
        return err
    }
    return nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
    _, err := ur.collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        log.Println("Error deleting user:", err)
        return err
    }
    return nil
}