package handlers

import (
	"net/http"
	"time"

	"go-transfer-agent/api-gateway/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

// --- request / response models ---

// CreateTransactionRequest is the JSON body for POST /fnd/transactions.
type CreateTransactionRequest struct {
	FromAccount string  `json:"from_account" binding:"required" example:"1234567890"`
	ToAccount   string  `json:"to_account" binding:"required" example:"0987654321"`
	Amount      float64 `json:"amount" binding:"required" example:"1500.50"`
	Currency    string  `json:"currency" binding:"required" example:"THB"`
	Memo        string  `json:"memo" example:"Payment for invoice #123"`
}

// TransactionResponse is the stub response for a created transaction.
type TransactionResponse struct {
	TransactionID string  `json:"transaction_id" example:"txn_abc123def456"`
	Status        string  `json:"status" example:"pending"`
	FromAccount   string  `json:"from_account" example:"1234567890"`
	ToAccount     string  `json:"to_account" example:"0987654321"`
	Amount        float64 `json:"amount" example:"1500.50"`
	Currency      string  `json:"currency" example:"THB"`
	CreatedBy     string  `json:"created_by" example:"admin"`
	CreatedAt     string  `json:"created_at" example:"2026-02-12T11:00:00Z"`
}

// --- handler ---

// CreateTransaction godoc
//
//	@Summary		Create fund transaction (stub)
//	@Description	Submits a new fund transfer transaction. This is a POC stub.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateTransactionRequest	true	"Transaction details"
//	@Success		201		{object}	TransactionResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/fnd/transactions [post]
func CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid request: " + err.Error()})
		return
	}

	// Extract the authenticated user from context.
	createdBy := "unknown"
	if raw, ok := c.Get("claims"); ok {
		if claims, ok := raw.(*middleware.UserClaims); ok {
			createdBy = claims.Sub
		}
	}

	// Stub response — no real persistence.
	c.JSON(http.StatusCreated, TransactionResponse{
		TransactionID: "txn_" + time.Now().Format("20060102150405"),
		Status:        "pending",
		FromAccount:   req.FromAccount,
		ToAccount:     req.ToAccount,
		Amount:        req.Amount,
		Currency:      req.Currency,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
}
