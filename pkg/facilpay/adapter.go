package facilpay

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

// FacilPayTransactionRecord matches Facil-Pay's merchant/escrow transaction export format.
type FacilPayTransactionRecord struct {
	TxUUID            string  `json:"tx_uuid"`
	MerchantKey       string  `json:"merchant_key"`
	CustomerPublicKey string  `json:"customer_public_key"`
	Amount            float64 `json:"amount"`
	Token             string  `json:"token"` // e.g. "USDC", "XLM"
	State             string  `json:"state"` // "SETTLED", "SETTLEMENT_PENDING", "FAILED"
	CreatedAt         string  `json:"created_at"`
	StellarMemo       string  `json:"stellar_memo"`
}

// FacilPayConnector adapts Facil-Pay transaction records.
type FacilPayConnector struct {
	SourcePath string
}

func NewFacilPayConnector(sourcePath string) *FacilPayConnector {
	return &FacilPayConnector{SourcePath: sourcePath}
}

func (f *FacilPayConnector) Name() string {
	return "facilpay"
}

func (f *FacilPayConnector) FetchInternalPayments(ctx context.Context, filter connector.Filter) ([]types.InternalPayment, error) {
	file, err := os.Open(f.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("facilpay connector error opening %s: %w", f.SourcePath, err)
	}
	defer file.Close()

	return f.parseRecords(file, filter)
}

func (f *FacilPayConnector) parseRecords(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	var records []FacilPayTransactionRecord
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, fmt.Errorf("facilpay connector error decoding json: %w", err)
	}

	var payments []types.InternalPayment
	for _, rec := range records {
		ts, err := time.Parse(time.RFC3339, rec.CreatedAt)
		if err != nil {
			ts = time.Now()
		}

		if !filter.TimeStart.IsZero() && ts.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && ts.After(filter.TimeEnd) {
			continue
		}

		p := types.InternalPayment{
			ID:          rec.TxUUID,
			SourceApp:   "facilpay",
			ReferenceID: rec.StellarMemo,
			Sender:      rec.CustomerPublicKey,
			Recipient:   rec.MerchantKey,
			Amount:      fmt.Sprintf("%.7f", rec.Amount),
			Asset:       rec.Token,
			Timestamp:   ts,
			Status:      rec.State,
		}

		payments = append(payments, p)
	}

	return payments, nil
}
