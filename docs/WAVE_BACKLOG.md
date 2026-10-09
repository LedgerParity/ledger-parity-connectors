# Focused engineering backlog

Six issue drafts reviewed against current source and existing open issues for the
October 9 submission. These are engineering tasks, not approval or earned points.
Proposed complexity is subject to maintainer review and Drips app configuration.

## 1. Reject duplicate JSON keys in canonical export decoding

## Description & Context

CLI input bounds do not automatically apply to library callers.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Define bounded shared decoding for library callers.
- [ ] Reject duplicate keys without changing decimal lexemes.
- [ ] Test nested metadata, valid empty arrays and oversized inputs.
- [ ] Coordinate with CLI consumers.
- [ ] Preserve current accepted field names.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/file/; pkg/database/; tests/

## Verification

`go test ./...; go vet ./...; go build ./...`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 2. Specify database fetch completeness and cancellation metadata

## Description & Context

A caller-supplied row function cannot establish completeness from row count alone.

## Proposed Complexity

High (200 points proposed); planning label `complexity: high`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Define explicit fetch coverage and limit/cancellation outcomes.
- [ ] Distinguish complete empty exports from truncated results.
- [ ] Reject contradictory scope/status filters.
- [ ] Preserve the existing API through a migration note.
- [ ] Do not bundle unreviewed live SQL credentials/drivers.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/database/; pkg/connector/; tests/

## Verification

`go test ./...; go vet ./...; go build ./...`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 3. Validate a sanitized SDP 7.0.0 export independently

## Description & Context

Current SDP adapter evidence is synthetic contract coverage.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Obtain a permitted sanitized export with pinned version/provenance.
- [ ] Independently establish sender/network/settlement scope.
- [ ] Compare field and timestamp semantics.
- [ ] Add negative drift fixtures.
- [ ] Preserve the unsupported verdict when evidence is absent.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/sdp/; docs/SDP.md; fixtures/

## Verification

`go test ./...; go vet ./...; go build ./...`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 4. Audit one named adapter against a pinned upstream contract

## Description & Context

Named wrappers must not imply unverified integration compatibility.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Choose one adapter and cite official schema/API revision.
- [ ] Map every economic/identity/state field with locators.
- [ ] Add permitted fixtures and drift rejection.
- [ ] Record unsupported cases.
- [ ] Keep live-service integration outside this scoped audit.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/stellopay/; pkg/trustlesswork/; pkg/facilpay/; docs/

## Verification

`go test ./...; go vet ./...; go build ./...`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 5. Test interval filters at export-window boundaries

## Description & Context

Explicit settlement intervals must remain conservative at filtered-window boundaries.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Identify current boundary coverage before adding tests.
- [ ] Include crossing, touching and excluded intervals for each supported representation.
- [ ] Ensure invalid rows cannot vanish through filtering.
- [ ] Preserve business references separately from transaction hashes.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/file/; pkg/sdp/; tests/

## Verification

`go test ./...; go vet ./...; go build ./...`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 6. Document canonical consumer error and precision examples

## Description & Context

Consumers need executable examples of exact amounts and rejected unsupported records.

## Proposed Complexity

Trivial (100 points proposed); planning label `complexity: trivial`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Add a compilable canonical JSON/CSV consumer example.
- [ ] Demonstrate decimal-string amounts and typed/error handling.
- [ ] Include issuer/network requirements and unsupported contract-token behavior.
- [ ] Run the example from a standalone checkout.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

README.md; examples/; pkg/connector/

## Verification

`go test ./...; go vet ./...; go build ./...`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## Published issue links

- [Reject duplicate JSON keys in canonical export decoding](https://github.com/LedgerParity/ledger-parity-connectors/issues/3)
- [Specify database fetch completeness and cancellation metadata](https://github.com/LedgerParity/ledger-parity-connectors/issues/4)
- [Validate a sanitized SDP 7.0.0 export independently](https://github.com/LedgerParity/ledger-parity-connectors/issues/5)
- [Audit one named adapter against a pinned upstream contract](https://github.com/LedgerParity/ledger-parity-connectors/issues/6)
- [Test interval filters at export-window boundaries](https://github.com/LedgerParity/ledger-parity-connectors/issues/7)
- [Document canonical consumer error and precision examples](https://github.com/LedgerParity/ledger-parity-connectors/issues/8)

These issues are published but have not been enrolled into Drips Wave or assigned.
