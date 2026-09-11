# Contributor tasks

1. Verify one named adapter's upstream contract. Deliver a revision-pinned field mapping and permitted sanitized fixture, or an explicit unsupported verdict. Acceptance: all economic/identity/state fields have locators, no invented fallbacks, tests reject contract drift. Do not add a live service until its schema and read-only authentication are reviewed.
2. Database export completeness. Deliver an explicit fetch-result coverage envelope and tests for caller truncation, cancellation and status-filter mismatches. Acceptance: partial rows cannot support a complete-export assertion. Keep SQL driver work separate.
3. Input size and duplicate JSON keys. Bound memory use and reject duplicate object keys without changing decimal lexemes. Acceptance: large/duplicate-key fixtures fail deterministically; empty valid arrays and nested metadata still work. Coordinate the shared decoder with CLI.

Ask the maintainer to confirm scope/availability before assignment; open a small PR with regression evidence. No Wave points or review-time promises are attached.
