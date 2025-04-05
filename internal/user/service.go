package user

import (
	"fmt"

	"github.com/144LMS/bet_master/models"
	"golang.org/x/crypto/bcrypt"
)

/*
Слой взаимодействия Репозитория с БД
*/

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
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
		Role:      user.Role,
	}, nil
}

func (s *UserService) CreateUserService(user *models.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return s.repo.CreateUserRepository(user)
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

func (s *UserService) Autenticate(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmailRepository(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.IsBanned {
		return nil, fmt.Errorf("user is banned")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
