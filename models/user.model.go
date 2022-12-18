package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SignUpInput struct {
	Name            string    `json:"name" bson:"name" bindinng:"required"`
	Email           string    `json:"email" bson:"email" binding:"required"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	ConfirmPassword string    `json:"confirm_password" bson:"confirm_password,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type DBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	ConfirmPassword string             `json:"confirm_password,omitempty" bson:"confirm_password,omitempty"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreateAt        time.Time          `json:"create_at" bson:"create_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreateAt,
		UpdatedAt: user.UpdatedAt,
	}
}
