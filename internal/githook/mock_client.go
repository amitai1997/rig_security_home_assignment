package githook

import "context"

// MockClient is a simple in-memory implementation of Client for tests.
type MockClient struct{}

// ListOrgRepositories returns a fixed list of repositories.
func (MockClient) ListOrgRepositories(ctx context.Context, org string) ([]Repository, error) {
	return []Repository{{Name: "repo1"}, {Name: "repo2"}}, nil
}

// ListCollaborators returns fixed collaborators based on repository name.
func (MockClient) ListCollaborators(ctx context.Context, owner, repo string) ([]Collaborator, error) {
	switch repo {
	case "repo1":
		return []Collaborator{
			{Login: "alice", Permission: "admin"},
			{Login: "bob", Permission: "write"},
		}, nil
	case "repo2":
		return []Collaborator{
			{Login: "carol", Permission: "read"},
		}, nil
	default:
		return nil, nil
	}
}
