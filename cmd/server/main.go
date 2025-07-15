package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/example/rig-security-svc/api/proto/v1"
	"github.com/example/rig-security-svc/internal/config"
	"github.com/example/rig-security-svc/internal/githook"
	"github.com/example/rig-security-svc/internal/policy"
	"github.com/example/rig-security-svc/internal/service"
)

func main() {
	cfg := config.LoadFromEnv()
	_ = cfg

	client := githook.MockClient{}
	engine, err := policy.NewEngine("permission == 'admin'")
	if err != nil {
		log.Fatalf("failed to init policy engine: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterRepositoryServiceServer(server, service.NewRepositoryService(client, engine))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("server exited: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	server.GracefulStop()
}
