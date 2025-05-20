package wallet

import (
	"errors"

	"github.com/144LMS/bet_master/models"
)

type WalletService struct {
	r *WalletRepository
}

func NewWalletService(r *WalletRepository) *WalletService {
	return &WalletService{r: r}
}

func (s *WalletService) GetWallet(userID uint) (*models.Wallet, error) {
	return s.r.GetWallet(userID)
}

func (s *WalletService) Deposit(userID uint, amount float64) (*models.Wallet, *models.Transaction, error) {
	if amount <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}
	return s.r.Deposit(userID, amount)
}

func (s *WalletService) Withdraw(userID uint, amount float64) (*models.Wallet, *models.Transaction, error) {
	if amount <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}
	return s.r.Withdraw(userID, amount)
}

func (s *WalletService) GetTransactions(userID uint) ([]models.Transaction, error) {
	return s.r.GetTransactions(userID)
}

func (s *WalletService) GetBalance(userID uint) (float64, error) {
	wallet, err := s.r.GetWallet(userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}
