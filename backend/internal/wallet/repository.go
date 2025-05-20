package wallet

import (
	"errors"
	//"fmt"

	"github.com/144LMS/bet_master/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetWallet(userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Preload("Transactions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&wallet, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Deposit(userID uint, amount float64) (*models.Wallet, *models.Transaction, error) {
	var wallet models.Wallet
	var tx *models.Transaction

	err := r.db.Transaction(func(txDB *gorm.DB) error {
		if err := txDB.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		if err := r.db.Save(&wallet).Error; err != nil {
			return err
		}

		tx = &models.Transaction{
			WalletID:    wallet.ID,
			Amount:      amount,
			Type:        "deposit",
			Description: "Пополнение баланса",
		}
		if err := r.db.Save(&tx).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return &wallet, tx, nil
}

func (r *WalletRepository) Withdraw(userID uint, amount float64) (*models.Wallet, *models.Transaction, error) {
	var wallet models.Wallet
	var tx *models.Transaction

	err := r.db.Transaction(func(txDb *gorm.DB) error {
		if err := txDb.Set("gorm:query_option", "FOR UPDATE").Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient funds")
		}

		wallet.Balance -= amount
		if err := txDb.Save(&wallet).Error; err != nil {
			return err
		}

		tx = &models.Transaction{
			WalletID:    wallet.ID,
			Amount:      amount,
			Type:        "withdrawal",
			Description: "Списание средств",
		}
		if err := txDb.Create(tx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &wallet, tx, nil
}

func (r *WalletRepository) GetTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	var wallet models.Wallet

	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}

	if err := r.db.Where("wallet_id = ?", wallet.ID).Order("created_at desc").Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
