package services

import (
	"otp-login-service/models"
	"otp-login-service/repository"
)

func GetUserByID(id int64) (*models.User, error) {
	return repository.FindByID(id)
}

func ListUsers(page, limit int, search string) ([]*models.User, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return repository.FindAll(page, limit, search)
}
func GetUserByPhoneNumber(phone string) (*models.User, error) {
	return repository.FindByPhoneNumber(phone)
}
