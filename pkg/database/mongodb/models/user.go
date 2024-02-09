package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Name      string             `json:"name,omitempty" bson:"name,omitempty"`
    Email     string             `json:"email,omitempty" bson:"email,omitempty"`
    Password  string             `json:"password,omitempty" bson:"password,omitempty"`
    AccessToken  string             `json:"access_token,omitempty" bson:"access_token,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
    CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
    UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}