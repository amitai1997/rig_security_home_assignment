# Rig Security Service

Rig Security Service is a minimal gRPC service written in Go. It lists the repositories in a GitHub organization and scans repository collaborators against a simple policy written in Google's [Common Expression Language](https://github.com/google/cel-go).

## Project Structure

Here's an overview of the important directories and files in this project:

- `api/proto/v1/rig.proto`: Defines the gRPC service contract, including the `RepositoryService` and the request/response messages (`ListRepositoriesRequest`, `ListRepositoriesResponse`, `RepositoryReport`, `PolicyViolation`). This is the primary interface for interacting with the service.
- `cmd/server/main.go`: The main application entry point. It initializes and configures the gRPC server, sets up the GitHub API client, and instantiates the policy engine.
- `internal/config/config.go`: Handles loading configuration from environment variables.
- `internal/githook/`: Contains the GitHub API client implementation.
    - `internal/githook/real_client.go`: Implements the `githook.Client` interface, providing methods to interact with the GitHub API (e.g., listing repositories and collaborators). It also includes rate limiting for GitHub API calls.
- `internal/policy/engine.go`: Implements the policy evaluation logic using Google's CEL. It defines the `Engine` interface and its `Scan` method, which evaluates a `Collaborator` against a defined CEL expression.
- `internal/service/repository.go`: Contains the core business logic for the `RepositoryService`. It orchestrates calls to the GitHub client to fetch repository and collaborator data, and then uses the policy engine to scan for violations.

## Request Flow

When a client makes a `ListRepositories` gRPC call to the service, the following general flow occurs:

1.  **Request Reception**: The `cmd/server/main.go` receives the gRPC request.
2.  **Service Invocation**: The request is routed to the `ListRepositories` method within `internal/service/repository.go`.
3.  **GitHub Data Fetching**: The `repository.go` service uses the `githook.Client` (implemented by `internal/githook/real_client.go`) to:
    *   List all repositories for the specified GitHub organization.
    *   For each repository, list its collaborators.
4.  **Concurrent Scanning**: To improve performance, the scanning of each repository's collaborators happens concurrently.
5.  **Policy Evaluation**: For each collaborator, the `internal/policy/engine.go` evaluates the collaborator's permissions against the configured CEL policy.
6.  **Violation Reporting**: If a policy violation is detected, it's included in the `RepositoryReport` for that repository.
7.  **Response**: The service compiles all `RepositoryReport` instances into a `ListRepositoriesResponse` and sends it back to the client.



Rig Security Service is a minimal gRPC service written in Go. It lists the repositories in a GitHub organization and scans repository collaborators against a simple policy written in Google's [Common Expression Language](https://github.com/google/cel-go).

## Prerequisites

- Go 1.22+
- Docker (for containerized execution)
- `protoc` with the Go plugins if you wish to regenerate protobuf code

## Configuration

Copy `.env.example` to `.env` and fill in the values:

```
GITHUB_TOKEN=your-token
GITHUB_ORG=github-organization
```

`GITHUB_TOKEN` must have permission to list repositories and collaborators in the organization.

## Local Setup

```bash
# build and run directly
GO111MODULE=on go run ./cmd/server
```

The server listens on `:50051`. Ensure the environment variables above are set before running.

## Docker Usage

Build and run using Docker Compose:

```bash
docker compose up --build
```

This uses the `Dockerfile` and `.env` file to create a small container image and start the gRPC server.

## API Usage

Call the service with `grpcurl`:

```bash
grpcurl -plaintext \
  -d '{"github_org":"example"}' \
  localhost:50051 \
  rig.v1.RepositoryService/ListRepositories
```

The response contains the repositories and any policy violations found.

## Policy Engine

Policies are simple CEL expressions. The default policy is defined in `cmd/server/main.go`:

```go
engine, err := policy.NewEngine("permission == 'admin'")
```

Modify this expression to change the rule. The `Collaborator` passed to the policy contains `login` and `permission` fields that you can reference in the expression.

