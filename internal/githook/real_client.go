package githook

import (
	"context"
	"net/http"

	github "github.com/google/go-github/v55/github"
	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

// rateLimitTransport wraps a RoundTripper with a rate limiter.
type rateLimitTransport struct {
	rt      http.RoundTripper
	limiter *rate.Limiter
}

func (t *rateLimitTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(r.Context()); err != nil {
		return nil, err
	}
	return t.rt.RoundTrip(r)
}

// RealClient implements Client using the GitHub API.
type RealClient struct {
	client *github.Client
}

// NewRealClient constructs a RealClient using the provided token and limiter.
func NewRealClient(token string, limiter *rate.Limiter) *RealClient {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.Background(), src)
	if limiter != nil {
		httpClient.Transport = &rateLimitTransport{rt: httpClient.Transport, limiter: limiter}
	}
	return &RealClient{client: github.NewClient(httpClient)}
}

// ListOrgRepositories returns repositories under the organization.
func (c *RealClient) ListOrgRepositories(ctx context.Context, org string) ([]Repository, error) {
	repos, _, err := c.client.Repositories.ListByOrg(ctx, org, nil)
	if err != nil {
		return nil, err
	}
	out := make([]Repository, 0, len(repos))
	for _, r := range repos {
		out = append(out, Repository{Name: r.GetName()})
	}
	return out, nil
}

// ListCollaborators returns collaborators for a repository.
func (c *RealClient) ListCollaborators(ctx context.Context, owner, repo string) ([]Collaborator, error) {
	users, _, err := c.client.Repositories.ListCollaborators(ctx, owner, repo, nil)
	if err != nil {
		return nil, err
	}
	collaborators := make([]Collaborator, 0, len(users))
	for _, u := range users {
		perm, _, err := c.client.Repositories.GetPermissionLevel(ctx, owner, repo, u.GetLogin())
		if err != nil {
			return nil, err
		}
		collaborators = append(collaborators, Collaborator{Login: u.GetLogin(), Permission: perm.GetPermission()})
	}
	return collaborators, nil
}
