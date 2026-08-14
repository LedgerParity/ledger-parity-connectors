package trustlesswork

import (
	"strings"
	"testing"
	"time"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
)

func TestTrustlessWorkAdapter_ParseRecords(t *testing.T) {
	rawJSON := `[
		{
			"escrow_id": "ESCROW-9001",
			"milestone_index": 1,
			"client_wallet": "GBCLIENTXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"contractor_wallet": "GBCONTRACTORXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"amount": 2500.50,
			"token_symbol": "USDC",
			"state": "RELEASED",
			"released_at": "2026-08-14T10:00:00Z",
			"transaction_hash": "tx_trustless_111",
			"contract_address": "CBTRUSTLESSWORKXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
		},
		{
			"escrow_id": "ESCROW-9001",
			"milestone_index": 2,
			"client_wallet": "GBCLIENTXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"contractor_wallet": "GBCONTRACTORXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			"amount": 1250.00,
			"token_symbol": "USDC",
			"state": "RELEASED",
			"released_at": "2026-08-14T11:30:00Z",
			"transaction_hash": "tx_trustless_222",
			"contract_address": "CBTRUSTLESSWORKXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
		}
	]`

	c := NewTrustlessWorkConnector("")
	payments, err := c.parseRecords(strings.NewReader(rawJSON), connector.Filter{})
	if err != nil {
		t.Fatalf("unexpected error parsing records: %v", err)
	}

	if len(payments) != 2 {
		t.Fatalf("expected 2 payments, got %d", len(payments))
	}

	p1 := payments[0]
	if p1.ID != "TW-ESCROW-9001-M1" {
		t.Errorf("expected ID TW-ESCROW-9001-M1, got %s", p1.ID)
	}
	if p1.Amount != "2500.5000000" {
		t.Errorf("expected Amount 2500.5000000, got %s", p1.Amount)
	}
	if p1.Asset != "USDC" {
		t.Errorf("expected Asset USDC, got %s", p1.Asset)
	}
	if p1.Metadata["escrow_id"] != "ESCROW-9001" {
		t.Errorf("expected metadata escrow_id ESCROW-9001, got %s", p1.Metadata["escrow_id"])
	}

	// Test with TimeFilter
	filter := connector.Filter{
		TimeStart: time.Date(2026, 8, 14, 10, 30, 0, 0, time.UTC),
		TimeEnd:   time.Date(2026, 8, 14, 12, 00, 0, 0, time.UTC),
	}
	filteredPayments, err := c.parseRecords(strings.NewReader(rawJSON), filter)
	if err != nil {
		t.Fatalf("unexpected filter error: %v", err)
	}
	if len(filteredPayments) != 1 {
		t.Fatalf("expected 1 payment after filter, got %d", len(filteredPayments))
	}
	if filteredPayments[0].ID != "TW-ESCROW-9001-M2" {
		t.Errorf("expected ID TW-ESCROW-9001-M2, got %s", filteredPayments[0].ID)
	}
}
