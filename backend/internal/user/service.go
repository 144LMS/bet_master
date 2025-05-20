package user

import (
	"fmt"

	"github.com/144LMS/bet_master/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) getUserService(id string) (*models.UserResponse, error) {
	user, err := s.repo.GetUserRepository(id)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        int(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
		Role:      user.Role,
	}, nil
}

func (s *UserService) getUserWithWalletService(id string) (*models.User, *models.Wallet, error) {
	user, err := s.repo.GetUserWithWalletRepository(id)
	if err != nil {
		return nil, nil, err
	}
	return user, &user.Wallet, nil
}

func (s *UserService) CreateUserWithWalletService(user *models.User) (*models.User, *models.Wallet, error) {
	if user.Password == "" {
		return nil, nil, fmt.Errorf("password cannot be empty")
	}

	var wallet *models.Wallet

	err := s.repo.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.CreateUserRepository(tx, user); err != nil {
			fmt.Println("Error creating user:", err)
			return err
		}

		createdWallet, err := s.repo.CreateWalletRepository(tx, user.ID)
		if err != nil {
			return err
		}

		wallet = createdWallet
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return user, wallet, nil
}

func (s *UserService) UpdateUserService(id string, req models.UpdateUserRequest) error {
	user, err := s.repo.GetUserRepository(id)
	if err != nil {
		return err
	}

	user.Username = req.Username
	user.Email = req.Email

	if req.Password != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashPassword)
	}
	return s.repo.UpdateUserRepository(user)
}

func (s *UserService) DeleteUserService(id string) error {
	return s.repo.DeleteUserRepository(id)
}

func (s *UserService) Authenticate(email, password string) (*models.User, error) {

	user, err := s.repo.GetUserByEmailRepository(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
