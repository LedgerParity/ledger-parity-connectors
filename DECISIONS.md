# Decisions

2026-09-11: Preserve existing MIT license/repository ownership. Canonical file input is the supported workflow. Named application packages remain experimental local-schema examples pending upstream verification (docs/PROVENANCE.md). Do not infer network, issuer or transaction identity from a product name. Decimal JSON numbers are preserved as json.Number; database floats are rejected; malformed dates never become current time. File limits fail rather than truncate silently. Core is pinned by real published Git revision, without sibling replacements.

CI now checks both the compatibility floor and stable Go, with manual dispatch and official updated action runtimes. This configuration is not a claim that remote CI ran. Generic database/example packages remain partial mapping helpers, while canonical file input provides strict validation.

Verification amendment: remote CI at 83fd7ec passed both Go matrix jobs, including the configured Linux race tests. Prior unverified-CI statements describe earlier inspection only. See PROJECT_HANDOFF.md for the run link and separate local evidence.
