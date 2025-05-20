package wallet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	s *WalletService
}

func NewWalletController(s *WalletService) *WalletController {
	return &WalletController{s: s}
}

func (c *WalletController) GetWallet(ctx *gin.Context) {
	userIDval, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	// Приводим к uint (или uint64, если так храните)
	userID, ok := userIDval.(uint)
	if !ok {
		// Если middleware кладёт uint64 — приведите к uint
		if id64, ok := userIDval.(uint64); ok {
			userID = uint(id64)
		} else if idInt, ok := userIDval.(int); ok {
			userID = uint(idInt)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID in context has wrong type"})
			return
		}
	}

	wallet, err := c.s.GetWallet(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           wallet.ID,
		"user_id":      wallet.UserID,
		"balance":      wallet.Balance,
		"created_at":   wallet.CreatedAt,
		"updated_at":   wallet.UpdatedAt,
		"transactions": wallet.Transactions,
	})
}

func (c *WalletController) Deposit(ctx *gin.Context) {
	userIDval, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	userID, ok := userIDval.(uint)
	if !ok {
		if id64, ok := userIDval.(uint64); ok {
			userID = uint(id64)
		} else if idInt, ok := userIDval.(int); ok {
			userID = uint(idInt)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID in context has wrong type"})
			return
		}
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, transaction, err := c.s.Deposit(userID, req.Amount)
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
	userIDval, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	userID, ok := userIDval.(uint)
	if !ok {
		if id64, ok := userIDval.(uint64); ok {
			userID = uint(id64)
		} else if idInt, ok := userIDval.(int); ok {
			userID = uint(idInt)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID in context has wrong type"})
			return
		}
	}

	var request struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, transaction, err := c.s.Withdraw(userID, request.Amount)
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
	userIDval, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	userID, ok := userIDval.(uint)
	if !ok {
		if id64, ok := userIDval.(uint64); ok {
			userID = uint(id64)
		} else if idInt, ok := userIDval.(int); ok {
			userID = uint(idInt)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID in context has wrong type"})
			return
		}
	}

	transactions, err := c.s.GetTransactions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (c *WalletController) GetBalance(ctx *gin.Context) {
	userIDval, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID not found in context"})
		return
	}
	userID, ok := userIDval.(uint)
	if !ok {
		if id64, ok := userIDval.(uint64); ok {
			userID = uint(id64)
		} else if idInt, ok := userIDval.(int); ok {
			userID = uint(idInt)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userID in context has wrong type"})
			return
		}
	}

	balance, err := c.s.GetBalance(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"balance": balance})
}
