# Connector provenance review

Reviewed 2026-09-11. Baseline adapters declare local field names but cite no versioned export schema, upstream source revision or captured response. Tests only establish parsing of synthetic local examples.

| Adapter | Primary source inspected | Finding | Claim allowed |
| --- | --- | --- | --- |
| Stellopay | https://github.com/Stellopay/stellopay-core | Describes Soroban payroll contracts; does not verify this adapter's payment_id/amount_xlm/processed_at export contract | Experimental local example only |
| Facil-Pay | https://github.com/Facil-Pay/facilpay-api | Real backend exists; inspected README does not establish tx_uuid/merchant_key/token export schema | Experimental local example only |
| Trustless Work | https://docs.trustlesswork.com/trustless-work/api-rest/introduction | Official API describes escrow and milestone operations; does not establish this adapter's flat milestone export format or fee/net amount interpretation | Experimental local example only |
| Canonical file | LedgerParity core types/validation and file tests | Project-owned JSON/CSV contract with exact amount and identity validation | Supported input contract for ordinary classic payments |
| Database rows | Caller-provided mapping/fetcher | No SQL driver or upstream DB schema is bundled; exact types enforced | Generic mapping helper, not a database product integration |

No exhaustive upstream code audit, authenticated API test, partnership or compatibility certification was performed. These source links establish where contract verification should begin, not that local fields match upstream. Soroban-based escrow releases cannot be reconciled by this preview's ordinary-payment ingestor. Legacy examples remain for migration; unsupported/incomplete records yield UNKNOWN in core, not matches.

To promote one adapter: pin an upstream schema/source revision; map every field including network, transaction/operation or event identity, exact amount units, asset issuer/contract, direction, fees and state; document exported-data completeness; add positive/negative fixtures with provenance and permission to redistribute. Validate against an upstream-provided sanitized example. Do not invent missing IDs or network.
