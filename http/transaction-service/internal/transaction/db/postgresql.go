package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"transaction_service/internal/transaction"
	"transaction_service/pkg/client/postgresql"
	"transaction_service/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r repository) Create(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error) {
	q := `INSERT INTO public.transaction (type, amount, comment) 
			VALUES ($1, $2, $3) 
			RETURNING id`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	queryRow := r.client.QueryRow(ctx, q, transaction.Type, transaction.Amount, transaction.Comment)
	if err := queryRow.Scan(&transaction.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return transaction, newErr
		}
		return transaction, err
	}

	return transaction, nil
}

func (r repository) FindAll(ctx context.Context) (t []transaction.Transaction, err error) {
	q := `
		SELECT id, type, amount, comment FROM public.transaction;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	transactions := make([]transaction.Transaction, 0)

	for rows.Next() {
		var t transaction.Transaction

		err = rows.Scan(&t.ID, &t.Type, &t.Amount, &t.Comment)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func (r repository) FindOne(ctx context.Context, id string) (transaction.Transaction, error) {
	q := `
		SELECT id, type, amount, comment FROM public.transaction WHERE id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var t transaction.Transaction
	queryRow := r.client.QueryRow(ctx, q, id)
	err := queryRow.Scan(&t.ID, &t.Type, &t.Amount, &t.Comment)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return t, nil
}

func NewRepository(client *pgxpool.Pool, logger *logging.Logger) transaction.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
