# Rig Security Service

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

