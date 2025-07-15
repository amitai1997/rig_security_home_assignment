package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/example/rig-security-svc/api/proto/v1"
	"github.com/example/rig-security-svc/internal/config"
	"github.com/example/rig-security-svc/internal/githook"
	"github.com/example/rig-security-svc/internal/policy"
	"github.com/example/rig-security-svc/internal/service"
	"golang.org/x/time/rate"
	"log/slog"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.LoadFromEnv()

	limiter := rate.NewLimiter(rate.Every(time.Second), 1)
	client := githook.NewRealClient(cfg.GitHubToken, limiter)

	engine, err := policy.NewEngine("permission == 'admin'")
	if err != nil {
		logger.Error("init engine", "error", err)
		return
	}

	server := grpc.NewServer()
	pb.RegisterRepositoryServiceServer(server, service.NewRepositoryService(client, engine))
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("listen", "error", err)
		return
	}

	go func() {
		logger.Info("server started", "addr", lis.Addr())
		if err := server.Serve(lis); err != nil {
			logger.Error("server exited", "error", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	server.GracefulStop()
}
