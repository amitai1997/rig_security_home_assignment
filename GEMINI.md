# GEMINI Rules and Guidelines


---

## api-design.mdc

---
description: 
globs: 
alwaysApply: true
---
**Endpoints & Resources**

* Design URIs around nouns, not verbs.
* Keep paths lowercase, kebab-case, plurals for collections.
* Nest resources only one level deep.
* Support filtering and pagination via query parameters.
* Avoid overloading endpoints with unrelated concerns.

**HTTP Semantics**

* Use standard verbs: GET, POST, PUT, PATCH, DELETE.
* Make GET safe and idempotent; PUT idempotent.
* Return proper status codes; prefer 4xx over 200 error payloads.
* Include Location header on resource creation.

**Payloads & Formats**

* Default to JSON; specify `application/json` Content-Type.
* Use consistent snake\_case or camelCase field names.
* Omit nulls; provide defaults when sensible.
* Document error envelope with code, message, details.
* Provide example requests/responses in docs.

**Versioning & Stability**

* Embed version in URI (`/v1/`) or header.
* Bump major only for breaking changes.
* Deprecate endpoints with sunset header and timeline.
* Maintain backward compatibility for at least one cycle.

**Security & Performance**

* Require TLS for all production traffic.
* Rate‑limit clients and expose quota headers.
* Compress responses with gzip or brotli.


---

## architecture-and-design.mdc

---
description: 
globs: 
alwaysApply: true
---
**Module Purpose & Cohesion**

* Give each module single, clear responsibility.
* Keep related functions and classes within same module.
* Limit exported symbols to the minimal public interface.
* Avoid circular dependencies; refactor shared logic upstream.
* Document module intent and invariants in top‑level docstring.

**File Length & Granularity**

* Limit source files to roughly 400 lines max.
* Split oversized files by responsibility boundaries.
* Extract inner classes to separate files when public.
* Co‑locate small, tightly coupled helpers with callers.
* Keep tests in matching file tree with concise files.

**Naming & Structure**

* Use descriptive, lowercase file names with underscores.
* Mirror package structure to domain language hierarchy.
* Append `_utils` only for generic helpers.
* Prefix experimental modules with `_draft` to flag instability.

**Architecture & Design**

* Modularize features into cohesive, loosely coupled components.
* Define shared abstractions early to prevent divergent implementations.
* Enforce single source of truth for configuration and constants.
* Prefer composition and parameterization over copy‑paste inheritance.
* Split domain layers (API, service, data, infra) into dedicated packages.
* Keep public APIs stable; hide internal helpers behind module‑private namespaces.
* Model business concepts explicitly; avoid anemic data‑only modules.
* Place cross‑cutting concerns (logging, metrics) in separate modules.
* Isolate third‑party integrations behind adapters to shield vendor changes.
* Favor functional purity in utility modules for easier reuse and testing.
* Establish clear module dependency direction; lower layers never import upward.
* Use interfaces or abstract base classes to invert dependencies when required.
* Validate layering rules with static analysis in CI.

**Dependency Management**

* Import only what you use; avoid wildcard imports.
* Prefer dependency injection over global imports.
* Group standard, third‑party, and local imports separately.
* Lock binary plugin versions to maintain ABI stability.
* Validate module graph with static analyzer in CI.

**Maintenance & Refactoring**

* Refactor duplicated logic into shared libraries promptly.
* Review file size and complexity metrics each sprint.
* Include deprecation notices in module docstring when moving APIs.
* Archive deprecated modules under `/legacy` with clear warnings.
* Update import paths after refactor using codemod tools.


---

## ci.mdc

---
description: 
globs: 
alwaysApply: true
---
**Pipeline Design**

* Trigger pipeline on every commit and pull request.
* Build in isolated, reproducible container images.
* Fail fast; abort pipeline after first critical error.
* Cache dependencies and artifacts to accelerate successive runs.

**Quality Gates**

* Run linters and formatters before executing tests.
* Require all unit and integration tests to pass.
* Enforce minimum code coverage threshold in pipeline.
* Scan dependencies and images for known vulnerabilities.

**Performance & Efficiency**

* Parallelize independent jobs to keep runtime under ten minutes.
* Use test‑impact analysis to execute only affected tests.

**Security & Compliance**

* Inject secrets at runtime from secure vault service.
* Sign build artifacts and generate checksums for verification.
* Record provenance metadata for every released artifact.

**Monitoring & Maintenance**

* Emit pipeline metrics and alert teams on failures.
* Auto‑rollback or block deployment on failed release stage.
* Keep CI configuration version‑controlled and peer‑reviewed.
* Review pipeline logs and flakiness trends each sprint.


---

## comment-and-docstrings.mdc

---
description: 
globs: 
alwaysApply: true
---
**Purpose & Clarity**

* Explain code intent and rationale, not obvious implementation.
* Use domain terminology consistently for clearer context.
* Avoid redundant comments that just restate code.

**Format & Style**

* Begin inline comment two spaces after statement.
* Capitalize first letter; end sentences with period.
* Separate block comments with preceding blank line.
* Use # TODO or # FIXME with initials and date.

**Docstring Structure**

* Enclose docstrings in triple double quotes.
* Start with concise summary line ending with period.
* Follow summary with blank line then details paragraph.
* List Args, Returns, Raises using Google style headers.
* Provide type hints inline with parameter descriptions.
* Include usage example or doctest when helpful.

**Coverage & Maintenance**

* Document every public module, class, function, method.
* Update comments and docstrings with every functional change.
* Delete obsolete comments during refactors.

**Tooling & Automation**

* Run pydocstyle, flake8-docstrings, or equivalent in CI.


---

## containerization-and-docker.mdc

---
description: 
globs: 
alwaysApply: true
---
**Image Build & Structure**

* Use multi‑stage builds to minimize final image size.
* Start from minimal, official base images.
* Pin base image digests for reproducible builds.
* Copy only needed files; maintain strict `.dockerignore`.
* Combine related `RUN` commands to reduce layers.


**Performance & Resource Management**

* Leverage build cache ordering for faster rebuilds.
* Add `HEALTHCHECK` instructions for runtime monitoring.
* Prefer distroless images for smaller attack surface.
* Use `tmpfs` mounts for ephemeral write paths.

**Configuration & Environment**

* Externalize configuration via environment variables only.
* Keep containers single‑process; one concern per image.
* Expose and document only required ports.
* Log to stdout/stderr; rely on platform aggregators.
* Avoid hard‑coded hostnames or IP addresses.

**CI/CD & Registry**

* Tag images semantically: version, latest, commit SHA.
* Promote images across environments using immutable tags.
* Prune unused images and dangling layers regularly.
* Automate rollback on failed deployment health checks.

**Docker Compose**

* Use Compose file latest stable version (3.9) for consistency.
* Reference images by digest or explicit tag, not `latest`.
* Define named volumes for persistent data; avoid anonymous volumes.
* Create dedicated networks to isolate service groups.
* Configure environment variables through a `.env` file; keep secrets out of compose.
* Add `healthcheck` sections for every service; rely on health status rather than `depends_on` alone.
* Set `restart` policies (`on-failure`, `unless-stopped`) appropriate to service role.
* Use profiles or multiple compose files to separate dev and prod overrides.
* Validate compose file in CI with `docker compose config`.
* Document each service purpose and exposed port in comments.
* Prefer build contexts pointing to Dockerfiles instead of inline build sections for reuse.
* Use named networks and explicit aliases for service discovery.


---

## databases-data-storage.mdc

---
description: 
globs: 
alwaysApply: false
---
**Data Modeling**

* Normalize relational schemas up to 3NF unless performance dictates denormalization.
* Model entities with stable natural keys; avoid business‑rule composite PKs.
* Use explicit ENUM or lookup tables for finite value sets.
* Embrace event sourcing or append‑only logs for audit requirements.
* Represent time zones with UTC timestamps; store timezone separately if needed.

**Storage Selection & Scalability**

* Choose relational for transactions, NoSQL for unstructured or high‑scale reads.
* Match consistency, availability, latency trade‑offs to business SLAs.
* Isolate analytical workloads on read replicas or data warehouse sinks.

**Schema & Migrations**

* Apply schema changes via version‑controlled migration scripts.
* Use declarative migration tools; avoid manual SQL in pipelines.
* Never drop columns without archival and deprecation phase.
* Run migrations in maintenance window or online rollout with blue/green.
* Validate migrations against staging snapshot before production.

**Performance & Indexing**

* Index columns used in WHERE, JOIN, ORDER BY clauses.
* Avoid over‑indexing; monitor write amplification and storage cost.
* Use covering indexes or materialized views for heavy read paths.
* Profile queries with EXPLAIN; refactor N+1 patterns.
* Cache hot reads with TTL layer; invalidate on write.


---

## directory-sturcture.mdc

---
description: 
globs: 
alwaysApply: true
---
**Structure**
* Start project with `src/`, `tests/`, `docs/`, and `infra/` directories.
* Place executable entry points inside `cmd/` or `bin/`.  
* Keep third-party patches under `vendor/`, outside `src/`.  
* Separate runtime data into `data/`, excluded from version control.  
* Use `/deploy` for Docker, scripts, and manifests.  

**Naming**
* Name directories lowercase, kebab-case or snake_case consistently.  
* Reserve `internal/` for non-public code.  
* Avoid generic names like `utils`; prefer descriptive domain terms.  
* Use plural nouns for collection folders, singular for leaf modules.  

**Modularity**
* Organize code by feature before technical layer.  
* Limit any directory to roughly fifteen files.  
* Prefix private submodules with `_` or use `.internal` suffix.  
* Extract reusable libraries into top-level `libs/`.  

**Testing & CI**
* Mirror `src/` hierarchy under `tests/` for one-to-one mapping.  
* Store test fixtures in `tests/fixtures/` separate from test code.  
* Keep CI pipelines in `.ci/` or `.github/workflows/`.  

**Documentation**
* Place `README.md`, `LICENSE`, `CHANGELOG.md` at repository root.  
* Add `CODEOWNERS` to clarify maintainership.  
* Provide folder-level READMEs for complex directories.  
* Track environment variables with `.env.example` template.  
"""
file_path = '/mnt/data/directory-sturcture-cursor-rules.md'
with open(file_path, 'w') as f:
    f.write(content)
file_path


---

## documentation-and-knowledge-sharing.mdc

---
description: 
globs: 
alwaysApply: true
---
**Scope & Intent**

* Document why, not just how.
* Target main audience; adjust depth accordingly.
* Keep docs close to code ownership boundaries.
* Version documentation with matching code releases.

**Writing Style**

* Use active voice and concise sentences.
* Prefer diagrams for complex flows over prose.
* Standardize terminology with project glossary.
* Include examples and expected outcomes for clarity.

**Artifacts & Formats**

* Use Markdown for code‑adjacent docs; adopt ADR template.
* Maintain API contracts in OpenAPI or GraphQL SDL.
* Store architecture diagrams as editable source files, not binaries.
* Provide README scaffolding for modules where needed.

**Governance & Automation**

* Archive outdated docs after deprecation window.


---

## git-branching-and-commits.mdc

---
description: 
globs: 
alwaysApply: true
---
**Branching Strategy**
* Use `main` as the production branch.
* Create feature branches from the latest `main`.
* Prefix branches with type, e.g., `feature/`, `bugfix/`, `hotfix/`, `release/`.
* Keep branches focused—one logical change per branch.
* Regularly pull latest `main` into feature branches to prevent drift.
* Delete merged branches to keep the repository clean.

**Commit Practices**
* Make atomic commits—each should represent one logical change.
* Write clear, imperative commit messages (e.g., “Fix login redirect bug”).
* Limit subject lines to 50 characters; add body if context is needed.
* Use present tense in commit messages ("Add", not "Added").
* Include issue/ticket number if applicable (e.g., `#123`).

**Rebasing & Merging**
* Use `rebase` for clean history before merging into `main`.
* Prefer `merge --no-ff` for preserving branch context in team projects.
* Always resolve conflicts before pushing rebased branches.

**Code Review & Collaboration**
* Open pull requests early for visibility and feedback.
* Ensure commits are linted and tested before opening a PR.
* Squash commits before merging if history clarity is desired.
* Use draft PRs for work-in-progress to signal early collaboration.

**Best Practices**
* Never commit directly to `main` or `production` branches.
* Don’t include unrelated changes in a single commit.
* Use `.gitignore` to exclude local config and secrets.
* Run `git status` before committing to double-check staged changes.
* Review commit history with `git log --oneline` to verify clarity.

---

## golang.mdc

---
description: 
globs: 
alwaysApply: true
---
Cursor Rules – Go Language

Syntax & Style
	•	Always run gofmt and goimports; commit formatted code only.
	•	Keep functions under 40 lines for readability.
	•	Name packages with short, lowercase, no underscores.
	•	Return early to reduce nesting; avoid deep indentation.
	•	Place variable declarations close to first use.
Concurrency: Goroutines & Channels
	•	Launch goroutines only when concurrent benefit outweighs overhead.
	•	Use context.Context for cancellation and deadlines across goroutines.
	•	Prefer unbuffered channels for sync; buffered for pipelines.
	•	Close channels from sender side after final write.
	•	Avoid sharing memory; communicate via channels for coordination.

Structs & Methods
	•	Use pointer receivers when modifying or copying large structs.
	•	Embed small structs to achieve composition over inheritance.
	•	Tag exported fields with JSON annotations for API clarity.
	•	Keep struct fields lowercase unless external access required.

Error Handling
	•	Return explicit errors; never silence or ignore them.
	•	Wrap errors with %w to preserve stack context.
	•	Avoid panics except for programmer errors or unrecoverable states.
	•	Provide sentinel errors for predictable branch decisions.

Performance & Memory
	•	Preallocate slices with make when size known.
	•	Use sync.Pool sparingly for short-lived object reuse.
	•	Benchmark critical paths using go test -bench.
	•	Profile CPU/heap with pprof before optimization.

Tooling & Build
	•	Manage dependencies with Go modules; commit go.sum.
	•	Run go vet and staticcheck in CI pipeline.
	•	Use go test -race to detect data races in concurrency.



---

## grpc.mdc

---
description: 
globs: 
alwaysApply: true
---
Cursor Rules – gRPC

API Definition & Protos
	•	Define services and messages in .proto files; keep them source-controlled.
	•	Use clear, versioned package namespaces (api.v1).
	•	Prefer explicit field numbers and reserve removed tags.
	•	Mark optional fields; avoid required to maintain compatibility.
	•	Document RPC methods with concise comments above declarations.

Service Design
	•	Keep RPC granular; avoid chatty method splitting.
	•	Favor unary calls; use streaming for large or real-time transfers.
	•	Model long-running jobs with server-side streaming progress updates.
	•	Set sensible deadlines on clients; propagate via context.

Error Handling
	•	Return standardized grpc.Status codes; map domain errors accordingly.
	•	Include machine-readable error_details for actionable responses.
	•	Avoid using UNKNOWN; choose most specific status code.

Performance & Streaming
	•	Compress payloads with grpc-encoding: gzip for large messages.
	•	Enable HTTP/2 keep-alive pings to detect dead peers.
	•	Chunk large streams; respect flow-control windows to prevent back-pressure.
	•	Reuse channel connections; maintain connection pools for high-QPS clients.

Tooling & CI
	•	Auto-generate client/server stubs during build; commit generated code only in vendoring cases.
	•	Run buf lint and prototool in CI to enforce style and breaking-change checks.
	•	Include golden proto examples and backward-compat tests per release.

---

## linting-formatting-and-type.checking.mdc

---
description: 
globs: 
alwaysApply: true
---
General Principles
	•	Adopt one canonical style guide per language and repository.
	•	Automate style enforcement; avoid subjective review debates.
	•	Run lightweight checks locally before committing.
	•	Document linter, formatter, type-checker usage in README.

Linting
	•	Enable linters with high severity thresholds; warn on stylistic issues.
	•	Fail CI on errors; block merges until resolved.
	•	Configure rule overrides centrally, version-controlled.
	•	Add custom plugins to capture domain-specific smells.

Formatting
	•	Use opinionated, zero-config formatters (Black, Prettier, gofmt).
	•	Reformat entire file on save or pre-commit.
	•	Keep line length consistent with tooling default (e.g., 88 for Python).
	•	Do not mix formatting changes with functional commits.

Type Checking
	•	Annotate public functions, classes, and exported symbols.
	•	Run static type checker in strict or incremental mode.
	•	Treat type errors as build-breaking defects.
	•	Prefer precise generics over Any; embrace gradual typing where needed.

Automation & CI
	•	Install tools via pinned versions in lockfiles.
	•	Integrate pre-commit hooks running lint, format, and type checks.
	•	Cache tool environments to speed CI pipelines.
	•	Generate badge displaying build and lint status in README.

---

## project-familiarity-and-avoiding-duplications.mdc

---
description: 
globs: 
alwaysApply: true
---
# Cursor Rules – Assistant Project Familiarity & Duplication Prevention

**Context Loading**

* Index entire codebase symbols, configs, and dependencies before any modification.
* Parse open PRs and feature branches to detect parallel changes.
* Ingest project documentation, ADRs, READMEs for domain context.
* Detect framework versions and coding conventions automatically.

**Impact Analysis**

* Resolve referenced identifiers to confirm existing definitions.
* Trace call graph to identify affected modules.
* Review test coverage gaps in impacted areas.
* Surface potential conflicting TODOs or technical debt.

**Code Generation & Reuse**

* Prefer extending existing abstractions over introducing new classes.
* Reuse helper utilities when functional overlap exceeds 70%.
* Consolidate duplicated logic into shared modules.
* Generate patches small and atomic per concern.

**Quality Gates**

* Run duplicate-code scanners on proposed diff pre-commit.
* Block commit if duplication metric exceeds configured threshold.
* Auto-run linter, formatter, and unit tests for changed files.
* Require developer approval for any new dependency introduction.

**Continuous Learning**

* Incrementally retrain code index on merged commits daily.
* Log prompt, rationale, and diff for auditability.
* Collect duplication trend metrics and suggest refactor tasks.


---

## testing.mdc

---
description: 
globs: 
alwaysApply: true
---
# Cursor Rules – Testing

**Strategy & Planning**

* Start testing during requirements to catch defects early.
* Use risk analysis to prioritize test scope.
* Keep a living test strategy updated each release.

**Unit Testing**

* Apply Arrange‑Act‑Assert for clear test structure.
* Isolate units by mocking external dependencies.
* Assert behavior, not implementation details, for resilience.
* Fail fast, deterministic tests run under any environment.

**Integration & End‑to‑End**

* Provision reproducible test environments using containers and IaC.
* Contract‑test service boundaries to detect breaking changes.
* Limit end‑to‑end suites to critical user journeys.
* Write e2e tests only where and when needed.
* Use masked or synthetic data for privacy compliance.

**Automation & CI/CD**

* Execute automated tests on every commit in CI.
* Parallelize suites to keep pipeline feedback under ten minutes.
* Tag and quarantine flaky tests immediately for repair.
* Collect coverage and trend metrics to guide improvement.

**Quality & Maintenance**

* Embed performance, load, and security checks in pipelines.
* Version‑control test code, data, and infrastructure together.
* Refactor test suites regularly to remove duplication.
* Review failing patterns to strengthen reliability and resilience.
* Update and expand existing tests whenever significant code or API changes are introduced to keep coverage accurate.
* Run the full test suite locally (or the relevant subset) after adding or modifying code to catch regressions before committing.

