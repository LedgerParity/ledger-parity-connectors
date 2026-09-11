package facilpay

import (
	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"strings"
	"testing"
)

func TestExactExampleAndBadTimestamp(t *testing.T) {
	raw := `[{"tx_uuid":"1","amount":922337203685.4775807,"created_at":"2026-09-01T12:00:00Z"}]`
	c := NewFacilPayConnector("")
	p, err := c.parseRecords(strings.NewReader(raw), connector.Filter{})
	if err != nil || len(p) != 1 || p[0].Amount != "922337203685.4775807" {
		t.Fatal(p, err)
	}
	for _, bad := range []string{strings.Replace(raw, "2026-09-01T12:00:00Z", "invalid", 1), strings.Replace(raw, "922337203685.4775807", "1.00000001", 1)} {
		if _, err := c.parseRecords(strings.NewReader(bad), connector.Filter{}); err == nil {
			t.Fatal("accepted invalid record")
		}
	}
}
