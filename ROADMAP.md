# Roadmap

Current: strict canonical JSON/CSV ingestion, exact decimal preservation, no fabricated timestamps, honest adapter provenance and standalone core dependency. Acceptance: canonical positive/negative tests, malformed/precision rejection, limit errors, go test/vet/build and isolated module consumption.

Next: independently verify one upstream contract and prove export completeness (docs/backlog.md). Soroban application adapters require a corresponding event reconciliation design before any integration claim. No additional repositories or signing workflow.

- [x] Canonical validation, amount/date fixes, provenance and bounded tasks implemented and pushed.
- [x] Isolated go test/vet/build using fresh module download passed.
- [ ] Upstream contract verification and operator adoption evidence.
- [x] Remote CI verified at the revisions linked in PROJECT_HANDOFF.md.
- [ ] Optional local Windows race execution (no C compiler); Linux CI race tests passed.
