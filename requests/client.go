package requests

import (
	"context"

	"github.com/i9si-sistemas/nine"
)

func NewWithContext(ctx context.Context) nine.Client {
	return nine.New(ctx)
}
