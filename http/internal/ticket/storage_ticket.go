package ticket

import "context"

type Repository interface {
	Create(ctx context.Context, ticket *Ticket) (*Ticket, error)
	FindAll(ctx context.Context) (t []Ticket, err error)
	FindOne(ctx context.Context, id string) (Ticket, error)
}
