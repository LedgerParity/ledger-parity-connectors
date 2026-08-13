# ledger-parity-connectors

[![CI](https://github.com/LedgerParity/ledger-parity-connectors/actions/workflows/ci.yml/badge.svg)](https://github.com/LedgerParity/ledger-parity-connectors/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`ledger-parity-connectors` is the adapter library for **LedgerParity**, providing target-app payment ingestors that transform proprietary database schemas and file exports into normalized `InternalPayment` models.

---

## 🔌 Supported Connectors & Adapters

- **Generic File Connector (`pkg/file`):** Ingests local CSV and JSON payment exports with flexible column mapping.
- **Stellopay Adapter (`pkg/stellopay`):** Ingests payment batch records from [Stellopay](https://github.com/search?q=stellopay) payroll platform exports.
- **Facil-Pay Adapter (`pkg/facilpay`):** Ingests merchant checkout and escrow transaction records from [Facil-Pay](https://github.com/search?q=facil-pay) payment exports.

---

## 📦 Usage Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/stellopay"
)

func main() {
	conn := stellopay.NewStellopayConnector("stellopay_export.json")
	payments, err := conn.FetchInternalPayments(context.Background(), connector.Filter{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fetched %d normalized internal payments\n", len(payments))
}
```

---

## 🧪 Testing

Run tests across all connector adapters:

```bash
go test -v ./...
```

---

## 🤝 Contributing

New application connectors (e.g. Betta-Pay, PayStell, Trustless-Work) can be added by implementing the `connector.Connector` interface. See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## 📄 License

[MIT License](LICENSE) © LedgerParity Maintainers.