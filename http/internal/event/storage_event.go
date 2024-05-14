package event

import "context"

type Repository interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	FindAll(ctx context.Context) (e []Event, err error)
	FindOne(ctx context.Context, id string) (Event, error)
	Update(ctx context.Context, event *Event, id string) error
	Delete(ctx context.Context, id string) error
}
