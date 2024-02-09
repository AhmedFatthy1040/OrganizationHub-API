package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationMember struct {
	Name        string `json:"name,omitempty" bson:"name,omitempty"`
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	AccessLevel string `json:"access_level,omitempty" bson:"access_level,omitempty"`
}

type Organization struct {
	ID           primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name,omitempty" bson:"name,omitempty"`
	Description  string               `json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt    time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Members      []OrganizationMember `json:"members,omitempty" bson:"members,omitempty"`
}
