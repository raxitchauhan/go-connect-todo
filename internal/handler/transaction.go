package handler

import (
	"context"
	"errors"
	"time"

	transactionv1 "go-connect-todo/gen/transaction/v1"
	"go-connect-todo/gen/transaction/v1/transactionv1connect"
	"go-connect-todo/internal/repository"
)

type transaction struct {
	repo repository.TransactionRepo
}

func NewTransaction(r repository.TransactionRepo) transactionv1connect.AccountServiceHandler {
	return &transaction{
		repo: r,
	}
}

// Deposit handles /deposit
func (t *transaction) Deposit(
	ctx context.Context,
	req *transactionv1.DepositRequest,
) (*transactionv1.DepositResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account_id is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if err := t.repo.Deposit(ctx, req.AccountId, req.Amount); err != nil {
		return nil, err
	}

	return &transactionv1.DepositResponse{
		Uuid: "some-uuid",
	}, nil
}

// Credit handles /credit
func (t *transaction) Credit(
	ctx context.Context,
	req *transactionv1.CreditRequest,
) (*transactionv1.CreditResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account_id is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if err := t.repo.Credit(ctx, req.AccountId, req.Amount); err != nil {
		return nil, err
	}

	return &transactionv1.CreditResponse{
		Uuid: "some-uuid",
	}, nil
}

// TransactionHistory handles /transaction_history
func (t *transaction) TransactionHistory(
	ctx context.Context,
	req *transactionv1.TransactionHistoryRequest,
) (*transactionv1.TransactionHistoryResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account_id is required")
	}
	if req.From == "" || req.To == "" {
		return nil, errors.New("from and to timestamps are required")
	}

	// parseTime and then send and handle errors later
	from, _ := time.Parse("", req.From)
	to, _ := time.Parse("", req.To)

	if to.Before(from) {
		return nil, errors.New("to must be after from")
	}

	txns, err := t.repo.GetTransactionHistory(
		ctx,
		req.AccountId,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	return &transactionv1.TransactionHistoryResponse{
		Transactions: txns,
	}, nil
}

// Balance handles /balance
func (t *transaction) Balance(
	ctx context.Context,
	req *transactionv1.BalanceRequest,
) (*transactionv1.BalanceResponse, error) {
	if req.AccountId == "" {
		return nil, errors.New("account_id is required")
	}

	balance, err := t.repo.GetBalance(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	return &transactionv1.BalanceResponse{
		Balance: balance,
	}, nil
}
