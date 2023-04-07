package nohup

import (
	"context"
	"os/signal"
	"syscall"
)

type IResident interface {
	Run(ctx context.Context, stopFn context.CancelFunc)
	Close()
}

func NewResident(parentCtx context.Context, programs ...IResident) {
	ctx, stop := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	for _, program := range programs {
		program.Run(ctx, stop)
	}
	<-ctx.Done()
	stop()

	for _, program := range programs {
		program.Close()
	}
}
