package database

import (
	"context"
	"testing"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
)

func TestDatabaseAdapter_MapRowToPayment(t *testing.T) {
	mapping := DefaultColumnMapping()

	mockFetcher := func(ctx context.Context, filter connector.Filter) ([]RowData, error) {
		return []RowData{
			{
				"id":                "DB-TX-1001",
				"sender_address":    "GBSENDERXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"recipient_address": "GBRECIPIENTXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"amount":            750.25,
				"currency":          "USDC",
				"created_at":        "2026-08-14T09:00:00Z",
				"status":            "CONFIRMED",
				"tx_hash":           "0xabc123mocktxhash",
			},
			{
				"id":                "DB-TX-1002",
				"sender_address":    "GBSENDERXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"recipient_address": "GBRECIPIENT2XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				"amount":            "100.00",
				"currency":          "XLM",
				"created_at":        "2026-08-14T10:00:00Z",
				"status":            "COMPLETED",
				"tx_hash":           "0xdef456mocktxhash",
			},
		}, nil
	}

	c := NewDatabaseConnector("postgres-payments", "app-db", mapping, mockFetcher)
	payments, err := c.FetchInternalPayments(context.Background(), connector.Filter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(payments) != 2 {
		t.Fatalf("expected 2 payments, got %d", len(payments))
	}

	p1 := payments[0]
	if p1.ID != "DB-TX-1001" {
		t.Errorf("expected ID DB-TX-1001, got %s", p1.ID)
	}
	if p1.Amount != "750.2500000" {
		t.Errorf("expected Amount 750.2500000, got %s", p1.Amount)
	}
	if p1.Asset != "USDC" {
		t.Errorf("expected Asset USDC, got %s", p1.Asset)
	}
	if p1.ReferenceID != "0xabc123mocktxhash" {
		t.Errorf("expected ReferenceID 0xabc123mocktxhash, got %s", p1.ReferenceID)
	}

	p2 := payments[1]
	if p2.Amount != "100.0000000" {
		t.Errorf("expected Amount 100.0000000, got %s", p2.Amount)
	}
	if p2.Asset != "XLM" {
		t.Errorf("expected Asset XLM, got %s", p2.Asset)
	}
}
