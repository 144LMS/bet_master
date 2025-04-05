package user

import (
	"github.com/144LMS/bet_master/models"
	"gorm.io/gorm"
)

/*
Слой взаимодействия с DB
*/

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) GetUserRepository(id string) (*models.User, error) {
	var user models.User

	err := repo.db.First(&user, id).Error
	return &user, err
}

func (repo *UserRepository) CreateUserRepository(user *models.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepository) UpdateUserRepository(user *models.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepository) DeleteUserRepository(id string) error {
	return repo.db.Delete(&models.User{}, id).Error
}

func (repo *UserRepository) GetUserByEmailRepository(email string) (*models.User, error) {
	var user models.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
