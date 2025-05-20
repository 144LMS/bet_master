package user

import (
	"errors"
	"fmt"
	"strconv"

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

	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format")
	}

	result := repo.db.Preload("Wallet").First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("database error: %v", result.Error)
	}

	if user.Wallet.ID == 0 {
		return nil, fmt.Errorf("user has no wallet")
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

func (r *UserRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
