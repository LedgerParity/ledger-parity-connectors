# Contributing to ledger-parity-connectors

Thank you for contributing adapter connectors to `ledger-parity-connectors`!

## How to add a new connector adapter

1. Create a new package under `pkg/<app_name>/`.
2. Define the source schema struct (matching your target app's JSON/CSV/DB format).
3. Implement the `connector.Connector` interface:
   ```go
   type Connector interface {
       Name() string
       FetchInternalPayments(ctx context.Context, filter Filter) ([]types.InternalPayment, error)
   }
   ```
4. Write unit tests validating schema parsing and filtering.
5. Open a Pull Request!
