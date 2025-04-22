package main

import (
	"common_library/logging"
	"common_library/metadata"
	"context"
	api2 "fileservice/pkg/api"
	"flag"
	"fmt"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"paymentservice/internal/config"
	"paymentservice/internal/data"
	"paymentservice/internal/db"
	"paymentservice/internal/handler"
	"paymentservice/internal/service"
	pb "paymentservice/pkg/api"
	api3 "schedule_service/pkg/api"
	"syscall"
	"userservice/pkg/api"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger := logging.New(zapLogger)

	ctx = logging.ContextWithLogger(ctx, logger)

	cfg, err := config.New()
	if err != nil {
		logger.Fatal(ctx, "cannot create config", zap.Error(err))
	}

	database, err := db.New(ctx, cfg)
	if err != nil {
		logger.Fatal(ctx, "cannot create db", zap.Error(err))
	}

	paymentRepo := data.NewPaymentRepository(database)
	stringPort := os.Getenv("GRPC_PORT")
	logger.Info(ctx, stringPort)
	serverAddr := flag.String("server", "localhost:"+stringPort, "gRPC server address")
	flag.Parse()
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(ctx, "failed to create grpc.NewClient"+err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			logger.Fatal(ctx, "failed tot close conn")
		}
	}(conn)

	userConn, err := grpc.NewClient("user-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	fileConn, err := grpc.NewClient("file-service:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	scheduleConn, err := grpc.NewClient("schedule-service:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))

	userClient := api.NewUserServiceClient(userConn)
	fileClient := api2.NewFileServiceClient(fileConn)
	scheduleClient := api3.NewScheduleServiceClient(scheduleConn)

	paymentService := service.NewPaymentService(paymentRepo, userClient, fileClient, scheduleClient)

	paymentHandler := handler.NewPaymentServiceServer(paymentService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal(ctx, "cannot create listener", zap.Error(err))
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			metadata.NewMetadataUnaryInterceptor(),
			logging.NewUnaryLoggingInterceptor(logger),
		)),
	)
	reflection.Register(server)
	pb.RegisterPaymentServiceServer(server, paymentHandler)

	logger.Info(ctx, "Starting gRPC server...", zap.Int("port", cfg.GRPCPort))
	go func() {
		if err := server.Serve(listener); err != nil {
			logger.Fatal(ctx, "failed to serve", zap.Error(err))
		}
	}()
	select {
	case <-ctx.Done():
		server.Stop()
		logger.Info(ctx, "Server Stopped")
	}
}
