package transaction

import "context"

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) (*Transaction, error)
	FindAll(ctx context.Context) (b []Transaction, err error)
	FindOne(ctx context.Context, id string) (Transaction, error)
}
