package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invitation struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	OrganizationID  primitive.ObjectID `json:"organization_id,omitempty" bson:"organization_id,omitempty"`
	InvitedEmail    string             `json:"invited_email,omitempty" bson:"invited_email,omitempty"`
	InvitationToken string             `json:"invitation_token,omitempty" bson:"invitation_token,omitempty"`
	ExpiresAt       time.Time          `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}