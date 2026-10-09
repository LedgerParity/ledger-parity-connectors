# Validated application exports for Stellar reconciliation

Prepared for the October 9, 2026 Stellar Wave submission.

## Purpose and implemented utility

The Go library normalizes canonical JSON/CSV and a version-scoped SDP export into exact payment records for ledger-parity-core. It validates input and preserves source/network/time provenance. The CLI consumes the pinned module.

## Reproduce the implementation

Go 1.22.2+; run `go test ./...`, `go vet ./...`, `go build ./...`. Follow README.md and docs/SDP.md for canonical input and the pinned SDP field contract.

## Evidence and supported scope

Local tests, vet, and build passed October 7. Named adapters include synthetic contract fixtures, not demonstrated provider integrations. Database integration uses a caller-supplied fetch function rather than a bundled SQL driver.

Evidence reference: [https://github.com/LedgerParity/ledger-parity-cli](https://github.com/LedgerParity/ledger-parity-cli).
Baseline source revision: `0e1b928937f41c24e43c40837f59348d410badc0`. Final reviewed preparation revision and CI
results belong in [VERIFICATION_OCT09.md](VERIFICATION_OCT09.md).

## Maintainers and contributor work

Maintainer: EthTobi; contact via GitHub, available anytime.
See [MAINTAINERS.md](../MAINTAINERS.md), [CONTRIBUTING.md](../CONTRIBUTING.md),
[SECURITY.md](../SECURITY.md), and [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md).
The [focused engineering backlog](WAVE_BACKLOG.md) describes real work, relevant
files, tests, and acceptance criteria. Draft complexity values require maintainer
review and app enrollment; they do not establish approval or earned points.

## Before applying

- Confirm the Drips Wave App covers this repository and check application slots.
- Publish the reviewed backlog issues and preserve links to their acceptance checks.
- Publish these preparation changes through a reviewed PR with passing CI.
- Apply under the implemented scope above; no production/adoption claims are implied.
