package server

import (
	"common_library/logging"
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"payment_service/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger := logging.New(zapLogger)

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

}
