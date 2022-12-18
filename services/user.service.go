package services

import "github.com/SarathLUN/auth-service-grpc-golang/models"

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
}
