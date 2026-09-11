package file

import (
	"context"
	"encoding/json"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func payment() types.InternalPayment {
	return types.InternalPayment{ID: "1", Network: "test", OperationType: "payment", Sender: "A", Recipient: "B", Amount: "922337203685.4775807", Asset: "XLM", AssetType: "native", Timestamp: time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC), Status: "completed"}
}
func TestJSONValidation(t *testing.T) {
	f := NewFileConnector("", "json", "test")
	p := payment()
	raw, _ := json.Marshal([]types.InternalPayment{p})
	rows, err := f.parseJSON(strings.NewReader(string(raw)), connector.Filter{})
	if err != nil || len(rows) != 1 || rows[0].Amount != p.Amount {
		t.Fatal(rows, err)
	}
	for _, bad := range []string{"null", string(raw) + " []", strings.Replace(string(raw), p.Amount, "1.00000001", 1), strings.Replace(string(raw), "2026-09-01T12:00:00Z", "invalid", 1), strings.Replace(string(raw), `"id":`, `"typo":`, 1)} {
		if _, err := f.parseJSON(strings.NewReader(bad), connector.Filter{}); err == nil {
			t.Fatal("accepted", bad)
		}
	}
	raw, _ = json.Marshal([]types.InternalPayment{p, p})
	if _, err = f.parseJSON(strings.NewReader(string(raw)), connector.Filter{Limit: 1}); err == nil {
		t.Fatal("silently truncated")
	}
}
func TestCSVValidationAndFilters(t *testing.T) {
	f := NewFileConnector("", "csv", "test")
	raw := "id,network,operation_type,sender,recipient,amount,asset,asset_type,timestamp,status\n1,test,payment,A,B,0.0000001,XLM,native,2026-09-01T12:00:00Z,completed\n"
	rows, err := f.parseCSV(strings.NewReader(raw), connector.Filter{})
	if err != nil || len(rows) != 1 || rows[0].Amount != "0.0000001" {
		t.Fatal(rows, err)
	}
	rows, err = f.parseCSV(strings.NewReader(raw), connector.Filter{Status: "pending"})
	if err != nil || len(rows) != 0 {
		t.Fatal(rows, err)
	}
	for _, bad := range []string{strings.Replace(raw, "timestamp", "id", 1), strings.Replace(raw, "2026-09-01T12:00:00Z", "", 1), strings.Replace(raw, "0.0000001", "NaN", 1), strings.Replace(raw, "asset_type", "typo", 1)} {
		if _, err = f.parseCSV(strings.NewReader(bad), connector.Filter{}); err == nil {
			t.Fatal("accepted", bad)
		}
	}
}
func TestFetchContextAndFormat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "input.json")
	os.WriteFile(path, []byte("[]"), 0600)
	f := NewFileConnector(path, "yaml", "test")
	if _, err := f.FetchInternalPayments(context.Background(), connector.Filter{}); err == nil {
		t.Fatal("unsupported format accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := f.FetchInternalPayments(ctx, connector.Filter{}); err == nil {
		t.Fatal("ignored cancellation")
	}
}

func TestIntervalCSVAndPartialFilter(t *testing.T) {
	f := NewFileConnector("", "csv", "test")
	raw := "id,network,operation_type,sender,recipient,amount,asset,asset_type,settlement_start,settlement_end,status,business_reference\n1,test,payment,A,B,1,XLM,native,2026-09-01T11:00:00Z,2026-09-01T13:00:00Z,completed,invoice-1\n"
	rows, err := f.parseCSV(strings.NewReader(raw), connector.Filter{})
	if err != nil || len(rows) != 1 || !rows[0].Timestamp.IsZero() || rows[0].BusinessReference != "invoice-1" {
		t.Fatal(rows, err)
	}
	if _, err = f.parseCSV(strings.NewReader(raw), connector.Filter{TimeStart: payment().Timestamp}); err == nil {
		t.Fatal("partial interval silently filtered")
	}
	rows, err = f.parseCSV(strings.NewReader(raw), connector.Filter{TimeStart: payment().Timestamp.Add(2 * time.Hour)})
	if err != nil || len(rows) != 0 {
		t.Fatal("nonoverlapping filter", rows, err)
	}
}
