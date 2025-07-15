package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/example/rig-security-svc/api/proto/v1"
	"github.com/example/rig-security-svc/internal/githook"
	"github.com/example/rig-security-svc/internal/policy"
)

// mockClient mocks githook.Client.
type mockClient struct{ mock.Mock }

func (m *mockClient) ListOrgRepositories(ctx context.Context, org string) ([]githook.Repository, error) {
	args := m.Called(ctx, org)
	return args.Get(0).([]githook.Repository), args.Error(1)
}

func (m *mockClient) ListCollaborators(ctx context.Context, owner, repo string) ([]githook.Collaborator, error) {
	args := m.Called(ctx, owner, repo)
	return args.Get(0).([]githook.Collaborator), args.Error(1)
}

// mockEngine mocks policy.Engine.
type mockEngine struct{ mock.Mock }

func (m *mockEngine) Scan(ctx context.Context, c policy.Collaborator) (*policy.Violation, error) {
	args := m.Called(ctx, c)
	if v := args.Get(0); v != nil {
		return v.(*policy.Violation), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestListRepositories(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(c *mockClient, e *mockEngine)
		expectErr  bool
		expectResp *pb.ListRepositoriesResponse
	}{
		{
			name: "single violation",
			setupMocks: func(c *mockClient, e *mockEngine) {
				repos := []githook.Repository{{Name: "repo1"}}
				c.On("ListOrgRepositories", mock.Anything, "myorg").Return(repos, nil)
				collaborators := []githook.Collaborator{{Login: "alice", Permission: "admin"}}
				c.On("ListCollaborators", mock.Anything, "myorg", "repo1").Return(collaborators, nil)
				e.On("Scan", mock.Anything, policy.Collaborator{Login: "alice", Permission: "admin"}).Return(&policy.Violation{Username: "alice", Permission: "admin", Rule: "deny"}, nil)
			},
			expectResp: &pb.ListRepositoriesResponse{Repositories: []*pb.RepositoryReport{{Name: "repo1", Violations: []*pb.PolicyViolation{{Username: "alice", Permission: "admin", Rule: "deny"}}}}},
		},
		{
			name: "no violations",
			setupMocks: func(c *mockClient, e *mockEngine) {
				repos := []githook.Repository{{Name: "repo1"}}
				c.On("ListOrgRepositories", mock.Anything, "myorg").Return(repos, nil)
				collaborators := []githook.Collaborator{{Login: "bob", Permission: "write"}}
				c.On("ListCollaborators", mock.Anything, "myorg", "repo1").Return(collaborators, nil)
				e.On("Scan", mock.Anything, policy.Collaborator{Login: "bob", Permission: "write"}).Return(nil, nil)
			},
			expectResp: &pb.ListRepositoriesResponse{Repositories: []*pb.RepositoryReport{{Name: "repo1"}}},
		},
		{
			name: "client error",
			setupMocks: func(c *mockClient, e *mockEngine) {
				c.On("ListOrgRepositories", mock.Anything, "myorg").Return([]githook.Repository(nil), errors.New("boom"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &mockClient{}
			e := &mockEngine{}
			tt.setupMocks(c, e)

			svc := NewRepositoryService(c, e)
			resp, err := svc.ListRepositories(context.Background(), &pb.ListRepositoriesRequest{GithubOrg: "myorg"})

			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectResp, resp)
			c.AssertExpectations(t)
			e.AssertExpectations(t)
		})
	}
}
