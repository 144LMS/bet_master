package bet

import (
	"errors"
	"net/http"

	"github.com/144LMS/bet_master/internal/wallet"
	"github.com/144LMS/bet_master/models"
	"github.com/gin-gonic/gin"
)

type BetController struct {
	betService    BetService
	walletService wallet.WalletService
}

func NewBetController(bs BetService, ws wallet.WalletService) *BetController {
	return &BetController{
		betService:    bs,
		walletService: ws,
	}
}

func (c *BetController) PlaceBet(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request struct {
		Amount   float64 `json:"amount" binding:"required,gt=0"`
		Odds     float64 `json:"odds" binding:"required,gt=1"`
		WalletID uint    `json:"wallet_id" binding:"required"`
		MatchID  uint    `json:"match_id" binding:"required"`
		Team     string  `json:"team" binding:"required,oneof=team1 team2 draw"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	match, err := c.betService.betRepo.GetMatch(request.MatchID)
	if err != nil || match.Status != models.MatchUpcoming {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or unavailable match"})
		return
	}

	wallet, err := c.walletService.GetWallet(request.WalletID)
	if err != nil || wallet.UserID != userID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet"})
		return
	}

	updatedWallet, withdrawTx, err := c.walletService.Withdraw(request.WalletID, request.Amount)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bet, err := c.betService.PlaceBet(models.Bet{
		UserID:       userID,
		WalletID:     request.WalletID,
		MatchID:      request.MatchID,
		Amount:       request.Amount,
		Odds:         request.Odds,
		PotentialWin: request.Amount * request.Odds,
		Team:         request.Team,
	})

	if err != nil {

		_, refundTx, refundErr := c.walletService.Deposit(request.WalletID, request.Amount)
		errorResponse := gin.H{
			"error":         err.Error(),
			"refund_status": refundErr == nil,
		}
		if refundTx != nil {
			errorResponse["refund_tx_id"] = refundTx.ID
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"bet":         bet,
		"match":       match,
		"wallet":      updatedWallet,
		"transaction": withdrawTx,
	})
}

func getUserIDFromContext(ctx *gin.Context) (uint, error) {
	const userIDKey = "userID"

	val, exists := ctx.Get(userIDKey)
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	switch v := val.(type) {
	case uint:
		return v, nil
	case int:
		if v < 0 {
			return 0, errors.New("invalid user ID: negative value")
		}
		return uint(v), nil
	case float64:
		if v < 0 {
			return 0, errors.New("invalid user ID: negative value")
		}
		return uint(v), nil
	default:
		return 0, errors.New("invalid user ID type in context")
	}
}
