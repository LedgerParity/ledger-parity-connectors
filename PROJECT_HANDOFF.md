# Project handoff

2026-09-11. Initial tree was clean at fe10911. Baseline tests passed but canonical file/facilpay had no tests. Implemented canonical file validation, preserved asset/network/operation identity, exact example amounts, database float rejection and malformed timestamp errors. Named adapters remain unverified local examples and may yield incomplete canonical records; CLI uses canonical file import instead.

Dependency: core 48e0036, fetched through Go tooling with real go.sum. No local replace or sibling import. Existing MIT license unchanged. See docs/PROVENANCE.md, docs/backlog.md and core docs/GAP_ASSESSMENT.md for sources and scope. User's recalled feedback concerns Stellar relevance/impact, not a supplied written rejection.

Verification in progress: local go test/vet/build; final results and remote evidence will be appended after execution. No live upstream adapter compatibility is claimed.

Local verification passed: go test ./..., go vet ./..., go build ./... on Go 1.24.4, using the downloaded pinned core module. Source changes are formatted with gofmt. Live upstream compatibility remains unverified. Remote CI results are not claimed.
