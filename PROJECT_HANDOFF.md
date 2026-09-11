# Project handoff

2026-09-11. Initial tree was clean at fe10911. Baseline tests passed but canonical file/facilpay had no tests. Implemented canonical file validation, preserved asset/network/operation identity, exact example amounts, database float rejection and malformed timestamp errors. Named adapters remain unverified local examples and may yield incomplete canonical records; CLI uses canonical file import instead.

Dependency: core 48e0036, fetched through Go tooling with real go.sum. No local replace or sibling import. Existing MIT license unchanged. See docs/PROVENANCE.md, docs/backlog.md and core docs/GAP_ASSESSMENT.md for sources and scope. User's recalled feedback concerns Stellar relevance/impact, not a supplied written rejection.

Verification in progress: local go test/vet/build; final results and remote evidence will be appended after execution. No live upstream adapter compatibility is claimed.

Local verification passed: go test ./..., go vet ./..., go build ./... on Go 1.24.4, using the downloaded pinned core module. Source changes are formatted with gofmt. Live upstream compatibility remains unverified. Remote CI results are not claimed.

Completion: local isolated archive test/vet/build passed on Go 1.24.4 with GOWORK=off and a fresh module cache (8 test cases; published core dependency downloaded). Implementation commit 902fa78 was pushed normally. The CLI also pins the newer core reconciliation changes directly through normal Go module version selection. No sibling checkout is needed.

CI now includes Go 1.22.2 and stable, manual dispatch and updated official action runtimes; remote success remains unverified. Local race execution is unavailable without a C compiler/CGO; the attempted portable compiler download was cancelled without installation. Core/CLI handoffs record the full evidence distinction. Named upstream contract validation and operator export completeness remain next tasks, not accomplished integrations.
