# Project handoff

Updated 2026-09-11. Canonical export milestone is implemented and verified. Initial tree was clean at fe10911. Code commit 902fa78 and CI/docs commit 83fd7ec were pushed normally. Existing MIT license is unchanged.

Canonical JSON/CSV input validates exact amounts, network/operation/asset identity and RFC3339 timestamps, rejects unknown fields/duplicate CSV columns and fails instead of silently truncating. Named example adapters preserve JSON decimal numbers and reject malformed dates; database row mapping rejects floating monetary values. Named Stellopay/Facil-Pay/Trustless Work packages remain unverified local examples that may yield incomplete canonical records. CLI selects canonical file import, not these examples. See docs/PROVENANCE.md.

Local test/vet/build passed on Go 1.24.4. An exported source archive passed with GOWORK=off and a fresh module cache: 8 test cases. Core dependency 48e0036 was downloaded through Go tooling with real go.sum; no sibling replace/import. CLI pins the newer core reconciliation revision directly through normal Go module selection.

Remote CI passed at 83fd7ec: https://github.com/LedgerParity/ledger-parity-connectors/actions/runs/34653797659 (two Go matrix jobs: minimum 1.22.2 and stable, configured Linux race tests/vet/build). This completion documentation commit changes no runtime/CI. Local Windows race execution lacked a C compiler; remote Linux evidence is distinct. Core and CLI verification/run links are recorded in CLI docs/VERIFICATION.md.

Next: independently verify one upstream contract using revision-pinned field mappings and a permitted sanitized fixture, or record an unsupported verdict. Resolve caller export completeness before asserting complete coverage. See docs/backlog.md and ROADMAP.md. No live upstream adapter compatibility, partnership, adoption or Drips approval is claimed. Original rejection text/application remain absent; the user's recollection concerns Stellar relevance/impact.
