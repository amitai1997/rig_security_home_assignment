package service

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"log/slog"

	pb "github.com/example/rig-security-svc/api/proto/v1"
	"github.com/example/rig-security-svc/internal/githook"
	"github.com/example/rig-security-svc/internal/policy"
)

// RepositoryService implements the RepositoryService gRPC server.
type RepositoryService struct {
	pb.UnimplementedRepositoryServiceServer
	client githook.Client
	engine policy.Engine
}

// NewRepositoryService constructs a RepositoryService.
func NewRepositoryService(c githook.Client, e policy.Engine) *RepositoryService {
	return &RepositoryService{client: c, engine: e}
}

// ListRepositories returns repository reports for the requested organization.
// It fetches repositories and scans each collaborator concurrently.
func (s *RepositoryService) ListRepositories(ctx context.Context, req *pb.ListRepositoriesRequest) (*pb.ListRepositoriesResponse, error) {
	reqID := uuid.New().String()
	logger := slog.With("request_id", reqID, "org", req.GetGithubOrg())
	logger.Info("handling list repositories request")

	repos, err := s.client.ListOrgRepositories(ctx, req.GetGithubOrg())
	if err != nil {
		logger.Error("list repositories", "error", err)
		return nil, err
	}

	resp := &pb.ListRepositoriesResponse{Repositories: make([]*pb.RepositoryReport, len(repos))}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	for i, r := range repos {
		wg.Add(1)
		go func(idx int, repo githook.Repository) {
			defer wg.Done()
			sem <- struct{}{}
			resp.Repositories[idx] = s.scanRepository(ctx, req.GetGithubOrg(), repo, logger)
			<-sem
		}(i, r)
	}

	wg.Wait()
	logger.Info("completed request")
	return resp, nil
}

func (s *RepositoryService) scanRepository(ctx context.Context, org string, repo githook.Repository, logger *slog.Logger) *pb.RepositoryReport {
	report := &pb.RepositoryReport{Name: repo.Name}
	collaborators, err := s.client.ListCollaborators(ctx, org, repo.Name)
	if err != nil {
		logger.Error("list collaborators", "repo", repo.Name, "error", err)
		return report
	}
	for _, c := range collaborators {
		v, err := s.engine.Scan(ctx, policy.Collaborator{Login: c.Login, Permission: c.Permission})
		if err != nil || v == nil {
			continue
		}
		report.Violations = append(report.Violations, &pb.PolicyViolation{
			Username:   v.Username,
			Permission: v.Permission,
			Rule:       v.Rule,
		})
		logger.Info("policy violation", "repo", repo.Name, "user", v.Username, "permission", v.Permission)
	}
	return report
}
