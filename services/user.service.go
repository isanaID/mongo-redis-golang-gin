package services

import "github.com/isanaID/mongo-redis-golang-gin/models"

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
}
