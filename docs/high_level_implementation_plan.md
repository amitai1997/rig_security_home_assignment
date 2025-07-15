Of course. Here is a comprehensive and detailed implementation plan for the Go home-assignment repository. This plan is designed to guide an AI agent through a five-day development sprint, ensuring adherence to best practices and the specific requirements of the task.

## **Project Implementation Plan: Rig Security gRPC Service**

### **Introduction & Guiding Principles**

This document outlines a detailed, five-stage implementation plan for the "Rig Security" home assignment. The primary objective is to develop a robust, high-quality gRPC service in Go that interacts with the GitHub API to list repositories and scan their access controls against a defined policy.

This plan is engineered for a single developer working within a five-day sprint. The following principles are paramount and will guide all development decisions:

  * **Simplicity and Minimalism:** We will implement exactly the features required by the assignment—nothing more. The goal is a focused, clean solution, not a feature-rich, production-ready platform. This aligns with the 6-8 hour estimated timeline and avoids the "build your own system" trap identified in the task analysis.
  * **Adherence to Standards:** All code and configuration will strictly follow the provided `.cursor/rules/` guidelines. This ensures consistency in API design, architecture, testing, containerization, and Go-specific idioms.
  * **Testability and Decoupling:** The architecture will be designed around interfaces to ensure that business logic is decoupled from external dependencies (like the GitHub API). This is critical for creating fast, reliable unit tests and is a hallmark of senior-level development.
  * **Clarity and Maintainability:** The code will be well-documented, clearly structured, and easy for a reviewer to understand, build, and run.

### **Technology Stack & Justification**

The technology stack is chosen to be modern, performant, and idiomatic to the Go ecosystem, directly addressing the challenges outlined in the task analysis.

| Component | Technology | Justification |
| :--- | :--- | :--- |
| **Language** | **Go 1.22** | Go's built-in support for concurrency with goroutines and channels is ideal for the I/O-bound nature of this task (making numerous calls to the GitHub API). Its static typing, single-binary deployment, and excellent standard library make it a perfect fit for a lean, performant microservice. |
| **API Framework** | **gRPC** | Mandated by the assignment. gRPC provides a high-performance, strongly-typed, contract-first approach to API development using Protocol Buffers (Protobuf). This is ideal for service-to-service communication. We will adhere to the principles in `rules/grpc.mdc`. |
| **GitHub API Client** | **`google/go-github`** | This is the canonical, feature-complete, and well-maintained library for interacting with the GitHub API in Go. Using it avoids the need to write boilerplate HTTP request/response handling and authentication logic. |
| **Policy Engine** | **`google/cel-go`** | The task analysis correctly identifies that building a custom engine is a trap and that Regex is insufficient. Open Policy Agent (OPA) is powerful but may be overkill for this task. Google's Common Expression Language (CEL) is a perfect middle ground: it's an embeddable, lightweight, yet powerful expression language that can be integrated directly into the Go binary without external dependencies, aligning with our minimalism principle. |
| **Testing** | **Go `testing` + `stretchr/testify`** | The standard `testing` package provides the core framework. The `testify` suite adds powerful assertion libraries (`assert`, `require`) and a first-class mocking toolkit (`mock`), which is essential for unit testing our service layer by mocking the GitHub client interface. This aligns with the strategy in `rules/testing.mdc`. |
| **Containerization** | **Docker & Docker Compose** | As noted in the analysis, providing a containerized solution is a "production readiness" bonus. A multi-stage `Dockerfile` will create a minimal, secure final image. `docker-compose.yml` will provide a one-command setup for reviewers. This follows the guidelines in `rules/containerization-and-docker.mdc`. |
| **Continuous Integration** | **GitHub Actions** | Tightly integrated with the source code repository, GitHub Actions provides a simple and effective way to automate our quality gates (linting, formatting, testing) on every commit, as specified in `rules/ci.mdc`. |

-----

## **Stage 1: Project Foundation & API Contract (Day 1)**

**Goal:** To establish a clean, idiomatic, and compilable project skeleton. This stage is about creating a solid foundation by defining the directory structure and the gRPC API contract before any business logic is written.

  * **Tasks for AI Agent:**

    1.  **Initialize Go Module:** Run `go mod init github.com/<your-username>/rig-security-svc`.
    2.  **Create Directory Structure:** Following `rules/directory-sturcture.mdc`, create the following layout:
        ```
        /
        ├── api/proto/v1/
        ├── cmd/server/
        ├── internal/config/
        ├── internal/githook/
        ├── internal/policy/
        ├── internal/service/
        └──.github/workflows/
        ```
    3.  **Define gRPC API Contract (`api/proto/v1/rig.proto`):**
          * Define the service contract. Adhere to the naming and versioning conventions in `rules/grpc.mdc`.
          * Create a `RepositoryService` with a `ListRepositories` RPC method.
          * Define the request message `ListRepositoriesRequest`, taking a `github_org` string.
          * Define the response message `ListRepositoriesResponse`, containing a repeated list of `RepositoryReport` messages.
          * Define the `RepositoryReport` message, which includes repository details and a list of `PolicyViolation` messages.
          * Define the `PolicyViolation` message with fields for the username, the permission level, and a description of the rule that was violated.
    4.  **Generate Go Stubs:** Run the `protoc` compiler with the Go and gRPC plugins to generate the server stubs and data types from the `.proto` file.
    5.  **Implement Configuration Loading (`internal/config/config.go`):**
          * Create a `Config` struct to hold `GitHubToken` and `GitHubOrg`.
          * Implement a `LoadFromEnv()` function to populate this struct from environment variables.
    6.  **Create Server Entry Point (`cmd/server/main.go`):**
          * Implement a `main` function that:
              * Loads the configuration.
              * Creates a new gRPC server.
              * Registers the (currently empty) service implementation.
              * Includes graceful shutdown logic to handle `SIGINT` and `SIGTERM` signals, allowing the server to finish in-flight requests before exiting.

  * **Acceptance Criteria:**

      * The project directory structure is in place.
      * The Go gRPC code is successfully generated from the `.proto` file.
      * The application compiles without errors (`go build./...`).
      * The server can be started with `go run./cmd/server` and shut down gracefully with `Ctrl+C`.

-----

## **Stage 2: Core Logic with Mocked Dependencies (Day 2)**

**Goal:** To implement and thoroughly unit-test the core business logic in complete isolation from external services. This is achieved by programming against interfaces and using mock implementations.

  * **Tasks for AI Agent:**

    1.  **Define Interfaces for Decoupling:**
          * In `internal/githook/client.go`, define a `Client` interface with methods like `ListOrgRepositories(ctx, org)` and `ListCollaborators(ctx, owner, repo)`.
          * In `internal/policy/engine.go`, define a `Engine` interface with a method `Scan(ctx, collaborator) (Violation, error)`.
    2.  **Implement Mock GitHub Client (`internal/githook/mock_client.go`):**
          * Create a struct that implements the `githook.Client` interface.
          * Hardcode its methods to return a static, predictable list of repositories and collaborators. This data will be the basis for our unit tests.
    3.  **Implement Policy Engine (`internal/policy/engine.go`):**
          * Create a CEL-based implementation of the `policy.Engine` interface.
          * It should be initialized with a policy string (e.g., `"collaborator.login == 'admin' && collaborator.permissions.admin == true"`).
          * The `Scan` method will evaluate the incoming collaborator data against the compiled CEL program.
    4.  **Implement the gRPC Service (`internal/service/repository.go`):**
          * Create the `RepositoryService` struct. It will receive the `githook.Client` and `policy.Engine` interfaces via its constructor (dependency injection).
          * Implement the `ListRepositories` gRPC method. This method will orchestrate the core logic:
              * Call the GitHub client to get the list of repositories.
              * Use a pool of goroutines to concurrently fetch the collaborators for each repository, respecting the concurrency principles from `rules/golang.mdc`.
              * For each collaborator, pass the data to the policy engine's `Scan` method.
              * Aggregate any reported violations into the final gRPC response.
    5.  **Write Unit Tests (`internal/service/repository_test.go`):**
          * Using `stretchr/testify/mock`, create mock implementations of the `Client` and `Engine` interfaces.
          * Write table-driven tests to cover various scenarios: no violations, single violation, multiple violations, errors from dependencies.
          * Assert that the service calls its dependencies with the correct arguments and correctly transforms the data into the final gRPC response.

  * **Acceptance Criteria:**

      * The `RepositoryService` is fully implemented.
      * Unit tests for the service achieve \>90% coverage and pass successfully.
      * The main server application, when run with the *mocked* GitHub client, can be queried with `grpcurl` and returns the expected, hardcoded results.

-----

## **Stage 3: Real-World Integration & Robustness (Day 3)**

**Goal:** To replace the mock GitHub client with a real implementation that communicates with the GitHub API, handling authentication, rate limiting, and errors gracefully.

  * **Tasks for AI Agent:**

    1.  **Implement Real GitHub Client (`internal/githook/real_client.go`):**
          * Create a new struct that implements the `githook.Client` interface.
          * Use the `google/go-github` library to implement the methods.
          * The client must be initialized with an `http.Client` that is configured to use the `GITHUB_TOKEN` from the environment for authentication.
    2.  **Implement Concurrency and Rate Limiting:**
          * When fetching collaborators for multiple repositories, use a `sync.WaitGroup` and a fixed-size pool of goroutines to control the level of concurrency and avoid overwhelming the GitHub API.
          * Wrap the `http.Client` used by the `go-github` library with a rate limiter (e.g., from `golang.org/x/time/rate`) to prevent hitting GitHub's API rate limits.
    3.  **Implement Structured Logging:**
          * Integrate the standard `log/slog` library throughout the application.
          * Add structured logs (with key-value pairs) for important events: server startup, handling a request (with request ID), policy violations found, and errors encountered.
    4.  **Update Server Entry Point (`cmd/server/main.go`):**
          * Modify the `main` function to instantiate and inject the *real* GitHub client into the `RepositoryService`.
    5.  **Write Integration Tests:**
          * Create a new test file, e.g., `internal/service/integration_test.go`, using the `//go:build integration` build tag.
          * These tests will run the service against the actual GitHub API (ideally using a dedicated test organization and token). They will verify the end-to-end flow is correct.

  * **Acceptance Criteria:**

      * The service can successfully authenticate with and fetch data from the GitHub API.
      * The service correctly applies policies to real data and returns accurate results.
      * The application logs are structured and provide useful operational insight.
      * Integration tests pass when run explicitly.

-----

## **Stage 4: Containerization & Continuous Integration (Day 4)**

**Goal:** To package the service into a lightweight, secure Docker image and to automate all quality checks in a CI pipeline.

  * **Tasks for AI Agent:**

    1.  **Create a Multi-Stage `Dockerfile`:**
          * Following `rules/containerization-and-docker.mdc`, create a `Dockerfile`.
          * **Builder Stage:** Use a `golang` base image to build the Go binary, ensuring dependencies are cached.
          * **Final Stage:** Use a minimal base image like `gcr.io/distroless/static-debian11` or `alpine`. Copy only the compiled Go binary and any necessary CA certificates from the builder stage. The target is a final image under 20MB.
    2.  **Create `docker-compose.yml`:**
          * Define a single service that builds the image from the `Dockerfile`.
          * Use `env_file` to load configuration from a local `.env` file, making it easy for reviewers to run.
    3.  **Set Up GitHub Actions CI (`.github/workflows/ci.yml`):**
          * Create a workflow that triggers on `push` and `pull_request` events.
          * Define a single job with the following sequential steps, adhering to `rules/ci.mdc`:
            1.  Check out code.
            2.  Set up Go.
            3.  Run `go vet./...` for static analysis.
            4.  Run `golangci-lint run` to enforce deeper linting rules.
            5.  Run `go test -v -race./...` to execute all unit tests with the race detector enabled.

  * **Acceptance Criteria:**

      * The command `docker compose up --build` successfully builds and starts the service.
      * The GitHub Actions pipeline is configured and passes on all quality gates.
      * A merge to `main` is blocked if any CI check fails.

-----

## **Stage 5: Documentation & Final Polish**

**Goal:** To produce comprehensive documentation that makes the project easy to understand, set up, and use, and to perform a final quality check.

  * **Tasks for AI Agent:**

    1.  **Write Comprehensive `README.md`:**
          * **Overview:** Briefly describe the project's purpose.
          * **Prerequisites:** List required tools (Go, Docker, `protoc`).
          * **Configuration:** Explain all environment variables, using `.env.example` as the template.
          * **Local Setup:** Provide step-by-step instructions for building and running the service locally.
          * **Docker Usage:** Explain how to run the service using `docker compose`.
          * **API Usage:** Provide clear `grpcurl` examples demonstrating how to call the `ListRepositories` RPC.
          * **Policy Engine:** Explain the basics of the CEL policy and where to modify it.
    2.  **Create `.env.example`:** Create a template file listing all required environment variables with placeholder values.
    3.  **Add GoDoc Comments:** Ensure all exported types, functions, and methods have clear, concise GoDoc comments explaining their purpose, parameters, and return values.
    4.  **Final Code Review:** Perform a final pass over the entire codebase to ensure all guidelines from the `.cursor/rules/` directory have been met, paying special attention to `golang.mdc` and `grpc.mdc`.

  * **Acceptance Criteria:**

      * The `README.md` is complete, accurate, and provides all necessary information for a reviewer.
      * The codebase is fully documented with GoDoc.
      * The project is clean, formatted, and ready for submission.