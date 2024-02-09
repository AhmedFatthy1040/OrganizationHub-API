package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

// SaveTokens saves the access token, refresh token, and their expiration time for a user.
func (ur *UserRepository) SaveTokens(ctx context.Context, userID primitive.ObjectID, accessToken, refreshToken string) error {
    // Define the filter to find the user by ID
    filter := bson.M{"_id": userID}

    // Define the update to set the access token, refresh token, and expiration time
    update := bson.M{
        "$set": bson.M{
            "access_token":  accessToken,
            "refresh_token": refreshToken,
        },
    }

    // Execute the update operation
    _, err := ur.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Println("Error saving tokens:", err)
        return err
    }

    return nil
}

// GetTokens retrieves the access token, refresh token, and their expiration time for a user.
func (ur *UserRepository) GetTokens(ctx context.Context, userID primitive.ObjectID) (string, string, error) {
    var user models.User
    projection := bson.M{"access_token": 1, "refresh_token": 1, "expires_at": 1}

    // Find the user by ID and retrieve the token fields
    err := ur.collection.FindOne(ctx, bson.M{"_id": userID}, options.FindOne().SetProjection(projection)).Decode(&user)
    if err != nil {
        log.Println("Error getting tokens:", err)
        return "", "", err
    }

    // Extract the token fields from the user document
    accessToken := user.AccessToken
    refreshToken := user.RefreshToken

    return accessToken, refreshToken, nil
}

// UpdateAccessToken updates the access token and its expiration time for a user.
func (ur *UserRepository) UpdateAccessToken(ctx context.Context, userID primitive.ObjectID, accessToken string, expiresAt time.Time) error {
    // Define the filter to find the user by ID
    filter := bson.M{"_id": userID}

    // Define the update to set the access token and its expiration time
    update := bson.M{
        "$set": bson.M{
            "access_token": accessToken,
            "expires_at":   expiresAt,
        },
    }

    // Execute the update operation
    _, err := ur.collection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Println("Error updating access token:", err)
        return err
    }

    return nil
}

func (ur *UserRepository) ValidateTokensByValue(ctx context.Context, refreshToken string) (primitive.ObjectID, error) {
    // Define the filter to find the user by refresh token
    filter := bson.M{"refresh_token": refreshToken}

    // Define the projection to only return the user's ID
    projection := bson.M{"_id": 1}

    // Execute the find one operation to get the user's ID by refresh token
    var result struct {
        UserID primitive.ObjectID `bson:"_id"`
    }

    err := ur.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
    if err != nil {
        return primitive.NilObjectID, err
    }

    return result.UserID, nil
}

// GetUserByToken retrieves a user by their access token from the database.
func (ur *UserRepository) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
    // Define the filter to find the user by access token
    filter := bson.M{"access_token": token}

    // Execute the query to find the user
    var user models.User
    err := ur.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        // Handle the error if the user is not found
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("user not found with token: %s", token)
        }
        return nil, err
    }

    return &user, nil
}
