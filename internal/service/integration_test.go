//go:build integration

package service

import (
	"context"
	"os"
	"testing"

	pb "github.com/example/rig-security-svc/api/proto/v1"
	"github.com/example/rig-security-svc/internal/githook"
	"github.com/example/rig-security-svc/internal/policy"
	"golang.org/x/time/rate"
)

func TestIntegrationListRepositories(t *testing.T) {
	org := os.Getenv("GITHUB_ORG")
	token := os.Getenv("GITHUB_TOKEN")
	if org == "" || token == "" {
		t.Skip("GITHUB_ORG or GITHUB_TOKEN not set")
	}
	client := githook.NewRealClient(token, rate.NewLimiter(rate.Every(rate.InfDuration), 1))
	engine, err := policy.NewEngine("permission == 'admin'")
	if err != nil {
		t.Fatal(err)
	}
	svc := NewRepositoryService(client, engine)
	_, err = svc.ListRepositories(context.Background(), &pb.ListRepositoriesRequest{GithubOrg: org})
	if err != nil {
		t.Fatalf("integration call failed: %v", err)
	}
}
