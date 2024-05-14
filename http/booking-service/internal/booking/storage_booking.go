package booking

import "context"

type Repository interface {
	Create(ctx context.Context, booking *Booking) (*Booking, error)
	FindAll(ctx context.Context) (b []Booking, err error)
	FindOne(ctx context.Context, id string) (Booking, error)
	Update(ctx context.Context, booking *Booking, id string) error
	Delete(ctx context.Context, id string) error
}
