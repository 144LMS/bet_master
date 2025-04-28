package user

import (
	"fmt"

	"github.com/144LMS/bet_master/models"
	"gorm.io/gorm"
)

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

func (repo *UserRepository) GetUserWithWalletRepository(id string) (*models.User, error) {
	var user models.User

	if err := repo.db.Preload("Wallet").First(&user, id); err != nil {
		return nil, fmt.Errorf("user with wallet not found: %v", err)
	}

	return &user, nil
}

func (repo *UserRepository) CreateUserRepository(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

func (repo *UserRepository) CreateWalletRepository(tx *gorm.DB, userID uint) (*models.Wallet, error) {
	wallet := models.Wallet{
		UserID:  userID,
		Balance: 0,
	}

	if err := tx.Create(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
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
