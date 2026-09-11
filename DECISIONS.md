# Decisions

2026-09-11: Preserve existing MIT license/repository ownership. Canonical file input is the supported workflow. Named application packages remain experimental local-schema examples pending upstream verification (docs/PROVENANCE.md). Do not infer network, issuer or transaction identity from a product name. Decimal JSON numbers are preserved as json.Number; database floats are rejected; malformed dates never become current time. File limits fail rather than truncate silently. Core is pinned by real published Git revision, without sibling replacements.
