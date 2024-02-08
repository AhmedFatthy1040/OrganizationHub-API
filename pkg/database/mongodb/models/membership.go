package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Membership struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID         primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	OrganizationID primitive.ObjectID `json:"organization_id,omitempty" bson:"organization_id,omitempty"`
	Role           string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt      time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}