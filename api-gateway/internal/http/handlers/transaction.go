package handlers

import (
	"context"
	"net/http"
	"time"

	"go-transfer-agent/api-gateway/internal/grpcclient"
	"go-transfer-agent/api-gateway/internal/http/middleware"
	fndv1 "go-transfer-agent/common/gen/fnd/v1"

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

// TransactionResponse is the REST response for a created transaction.
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

// TransactionListResponse contains a list of transactions.
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total" example:"1"`
}

// --- handler struct ---

// FNDHandler holds the gRPC client for FND service.
type FNDHandler struct {
	fnd *grpcclient.FNDClient
}

// NewFNDHandler creates a new FNDHandler with gRPC client dependency.
func NewFNDHandler(fnd *grpcclient.FNDClient) *FNDHandler {
	return &FNDHandler{fnd: fnd}
}

// --- handlers ---

// CreateTransaction godoc
//
//	@Summary		Create fund transaction
//	@Description	Submits a new fund transfer transaction via FND gRPC service.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateTransactionRequest	true	"Transaction details"
//	@Success		201		{object}	TransactionResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/fnd/transactions [post]
func (h *FNDHandler) CreateTransaction(c *gin.Context) {
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

	// Convert amount to minor units (satang: *100)
	amountMinor := int64(req.Amount * 100)

	// Call FND gRPC service
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.fnd.Client.CreateTransaction(ctx, &fndv1.CreateTransactionRequest{
		FromAccount: req.FromAccount,
		ToAccount:   req.ToAccount,
		AmountMinor: amountMinor,
		Currency:    req.Currency,
		Memo:        req.Memo,
		CreatedBy:   createdBy,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "FND service error: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, TransactionResponse{
		TransactionID: resp.GetTransactionId(),
		Status:        resp.GetStatus(),
		FromAccount:   req.FromAccount,
		ToAccount:     req.ToAccount,
		Amount:        req.Amount,
		Currency:      req.Currency,
		CreatedBy:     createdBy,
		CreatedAt:     resp.GetCreatedAt(),
	})
}

// ListTransactions godoc
//
//	@Summary		List fund transactions
//	@Description	Returns a list of fund transactions from FND gRPC service.
//	@Tags			transactions
//	@Produce		json
//	@Success		200	{object}	TransactionListResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/fnd/transactions [get]
func (h *FNDHandler) ListTransactions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.fnd.Client.ListTransactions(ctx, &fndv1.ListTransactionsRequest{
		PageSize: 50,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "FND service error: " + err.Error()})
		return
	}

	txns := make([]TransactionResponse, 0, len(resp.GetTransactions()))
	for _, t := range resp.GetTransactions() {
		txns = append(txns, TransactionResponse{
			TransactionID: t.GetTransactionId(),
			Status:        t.GetStatus(),
			FromAccount:   t.GetFromAccount(),
			ToAccount:     t.GetToAccount(),
			Amount:        float64(t.GetAmountMinor()) / 100.0,
			Currency:      t.GetCurrency(),
			CreatedBy:     t.GetCreatedBy(),
			CreatedAt:     t.GetCreatedAt(),
		})
	}

	c.JSON(http.StatusOK, TransactionListResponse{
		Transactions: txns,
		Total:        int(resp.GetTotal()),
	})
}
