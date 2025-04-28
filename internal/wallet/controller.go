package wallet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	s *WalletService
}

func NewWalletController(s *WalletService) *WalletController {
	return &WalletController{s: s}
}

func (c *WalletController) GetWallet(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	wallet, err := c.s.GetWallet(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
	}

	ctx.JSON(http.StatusOK, wallet)
}

func (c *WalletController) Deposit(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, transaction, err := c.s.Deposit(uint(userID), req.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"wallet":      wallet,
		"transaction": transaction,
	})
}

func (c *WalletController) Withdraw(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var request struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, transaction, err := c.s.Withdraw(uint(userID), request.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"wallet":      wallet,
		"transaction": transaction,
	})
}

func (c *WalletController) GetTransactions(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	transactions, err := c.s.GetTransactions(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (c *WalletController) GetBalance(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	balance, err := c.s.GetBalance(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}
