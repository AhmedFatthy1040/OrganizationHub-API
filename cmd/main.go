package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/handlers"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/middleware"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/api/routes"
	"github.com/AhmedFatthy1040/OrganizationHub-API/pkg/database/mongodb/repository"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	usersCollection     *mongo.Collection
	organizationsCollection *mongo.Collection
	tokensCollection *mongo.Collection
	invitationsCollection *mongo.Collection
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI not found in .env file")
	}

	// Connect to MongoDB
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")

	// Set database
	databaseName := os.Getenv("DATABASE_NAME")
	database = client.Database(databaseName)

	// Set collections
	usersCollectionName := os.Getenv("USERS_COLLECTION")
	organizationsCollectionName := os.Getenv("ORGANIZATIONS_COLLECTION")
	tokensCollectionName := os.Getenv("TOKENS_COLLECTION")
	invitationsCollectionName := os.Getenv("INVITATIONS_COLLECTION")

	usersCollection = database.Collection(usersCollectionName)
	organizationsCollection = database.Collection(organizationsCollectionName)
	tokensCollection = database.Collection(tokensCollectionName)
	invitationsCollection = database.Collection(invitationsCollectionName)
}

func main() {
    // Initialize Gin router
	router := gin.Default()

	// Initialize repositories
    userRepository := repository.NewUserRepository(database)
	organizationRepository := repository.NewOrganizationRepository(database)

	organizationHandler := handlers.NewOrganizationHandler(organizationRepository)

	// Setup middleware
    router.Use(middleware.BearerTokenAuth(userRepository, organizationRepository))

    // Setup routes
    routes.SetupUserRoutes(router, userRepository)
	routes.SetupOrganizationRoutes(router, organizationHandler)

	// Start the Gin server
	router.Run(":8080")
}