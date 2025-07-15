package githook

import "context"

// Repository represents a GitHub repository.
type Repository struct {
	Name string
}

// Collaborator represents a repository collaborator with a permission level.
type Collaborator struct {
	Login      string
	Permission string
}

// Client defines methods to interact with GitHub.
type Client interface {
	// ListOrgRepositories returns repositories under the organization.
	ListOrgRepositories(ctx context.Context, org string) ([]Repository, error)
	// ListCollaborators returns collaborators for the given repository.
	ListCollaborators(ctx context.Context, owner, repo string) ([]Collaborator, error)
}
