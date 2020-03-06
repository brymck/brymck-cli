package commands

import (
	"context"
	"time"
)

func makeContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}
