package repository

import (
	"context"
	"database/sql"
	"time"

	transactionv1 "go-connect-todo/gen/transaction/v1"

	"github.com/google/uuid"
)

type TransactionRepo interface {
	Deposit(ctx context.Context, accountID string, amount int64) error
	Credit(ctx context.Context, accountID string, amount int64) error
	GetTransactionHistory(
		ctx context.Context,
		accountID string,
		from time.Time,
		to time.Time,
	) ([]*transactionv1.Transaction, error)
	GetBalance(ctx context.Context, accountID string) (int64, error)
}

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Deposit(
	ctx context.Context,
	accountID string,
	amount int64,
) error {
	query := `
		INSERT INTO transactions (
			uuid,
			account_id,
			amount,
			is_credit
		) VALUES ($1, $2, $3, false)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		uuid.NewString(), // txn uuid
		accountID,
		amount,
	)

	return err
}

func (r *transactionRepo) Credit(
	ctx context.Context,
	accountID string,
	amount int64,
) error {
	query := `
		INSERT INTO transactions (
			uuid,
			account_id,
			amount,
			is_credit
		) VALUES ($1, $2, $3, true)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		uuid.NewString(), // txn uuid
		accountID,
		amount,
	)

	return err
}

func (r *transactionRepo) GetTransactionHistory(
	ctx context.Context,
	accountID string,
	from time.Time,
	to time.Time,
) ([]*transactionv1.Transaction, error) {
	query := `
		SELECT
			id,
			uuid,
			account_id,
			amount,
			created_at,
			is_credit
		FROM transactions
		WHERE account_id = $1
		  AND created_at >= $2
		  AND created_at <= $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, accountID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []*transactionv1.Transaction

	for rows.Next() {
		var (
			id        string
			uuidStr   string
			accID     string
			amount    int64
			createdAt time.Time
			isCredit  bool
		)

		if err := rows.Scan(
			&id,
			&uuidStr,
			&accID,
			&amount,
			&createdAt,
			&isCredit,
		); err != nil {
			return nil, err
		}

		txns = append(txns, &transactionv1.Transaction{
			Id:        id,
			Uuid:      uuidStr,
			AccountId: accID,
			Amount:    amount,
			CreatedAt: createdAt.Format(time.RFC3339),
			IsCredit:  isCredit,
		})
	}

	return txns, rows.Err()
}

func (r *transactionRepo) GetBalance(
	ctx context.Context,
	accountID string,
) (int64, error) {
	query := `
		SELECT
			COALESCE(
				SUM(
					CASE
						WHEN is_credit THEN amount
						ELSE -amount
					END
				),
				0
			) AS balance
		FROM transactions
		WHERE account_id = $1
	`

	var balance int64
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
