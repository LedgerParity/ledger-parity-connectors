package stellopay_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/stellopay"
)

func TestStellopayConnector(t *testing.T) {
	tempDir := t.TempDir()
	jsonFile := filepath.Join(tempDir, "stellopay.json")

	content := `[
		{
			"payment_id": "stell_1",
			"batch_id": "batch_99",
			"employee_wallet": "GEMP123",
			"employer_wallet": "GEMP456",
			"amount_xlm": 150.0,
			"currency": "XLM",
			"status": "COMPLETED",
			"processed_at": "2026-08-13T10:00:00Z",
			"stellar_tx_hash": "txhash123"
		}
	]`

	if err := os.WriteFile(jsonFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed creating test json: %v", err)
	}

	conn := stellopay.NewStellopayConnector(jsonFile)
	payments, err := conn.FetchInternalPayments(context.Background(), connector.Filter{})

	if err != nil {
		t.Fatalf("unexpected error fetching payments: %v", err)
	}
	if len(payments) != 1 {
		t.Fatalf("expected 1 payment record, got %d", len(payments))
	}

	p := payments[0]
	if p.ID != "stell_1" {
		t.Errorf("expected ID stell_1, got %s", p.ID)
	}
	if p.Recipient != "GEMP123" {
		t.Errorf("expected Recipient GEMP123, got %s", p.Recipient)
	}
	if p.SourceApp != "stellopay" {
		t.Errorf("expected SourceApp stellopay, got %s", p.SourceApp)
	}
}
