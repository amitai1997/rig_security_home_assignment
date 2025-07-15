package service

import (
	"context"
	"sync"

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

// ListRepositories implements the core logic.
func (s *RepositoryService) ListRepositories(ctx context.Context, req *pb.ListRepositoriesRequest) (*pb.ListRepositoriesResponse, error) {
	repos, err := s.client.ListOrgRepositories(ctx, req.GetGithubOrg())
	if err != nil {
		return nil, err
	}

	resp := &pb.ListRepositoriesResponse{}
	resp.Repositories = make([]*pb.RepositoryReport, len(repos))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, r := range repos {
		resp.Repositories[i] = &pb.RepositoryReport{Name: r.Name}
		wg.Add(1)
		go func(idx int, repo githook.Repository) {
			defer wg.Done()
			collaborators, err := s.client.ListCollaborators(ctx, req.GetGithubOrg(), repo.Name)
			if err != nil {
				return
			}
			for _, c := range collaborators {
				v, err := s.engine.Scan(ctx, policy.Collaborator{Login: c.Login, Permission: c.Permission})
				if err != nil || v == nil {
					continue
				}
				mu.Lock()
				resp.Repositories[idx].Violations = append(resp.Repositories[idx].Violations, &pb.PolicyViolation{
					Username:   v.Username,
					Permission: v.Permission,
					Rule:       v.Rule,
				})
				mu.Unlock()
			}
		}(i, r)
	}

	wg.Wait()
	return resp, nil
}
