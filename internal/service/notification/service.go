package notification

import (
	"context"
)

type Service interface {
}

type service struct {
	ctx context.Context
}

func NewService(ctx context.Context) (*service, error) {
	srv := service{
		ctx: ctx,
	}

	return &srv, nil
}
