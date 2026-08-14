package trustlesswork

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

// TrustlessWorkMilestoneRecord represents an escrow milestone release event in Trustless Work.
type TrustlessWorkMilestoneRecord struct {
	EscrowID         string  `json:"escrow_id"`
	MilestoneIndex   int     `json:"milestone_index"`
	ClientWallet     string  `json:"client_wallet"`
	ContractorWallet string  `json:"contractor_wallet"`
	Amount           float64 `json:"amount"`
	TokenSymbol      string  `json:"token_symbol"` // e.g. "USDC", "XLM", "EURC"
	State            string  `json:"state"`        // "RELEASED", "DISPUTED", "RESOLVED", "REFUNDED"
	ReleasedAt       string  `json:"released_at"`
	TransactionHash  string  `json:"transaction_hash"`
	ContractAddress  string  `json:"contract_address"`
}

// TrustlessWorkConnector adapts Trustless Work escrow milestone payment data.
type TrustlessWorkConnector struct {
	SourcePath string
}

func NewTrustlessWorkConnector(sourcePath string) *TrustlessWorkConnector {
	return &TrustlessWorkConnector{SourcePath: sourcePath}
}

func (t *TrustlessWorkConnector) Name() string {
	return "trustless-work"
}

func (t *TrustlessWorkConnector) FetchInternalPayments(ctx context.Context, filter connector.Filter) ([]types.InternalPayment, error) {
	file, err := os.Open(t.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("trustless-work connector failed opening file %s: %w", t.SourcePath, err)
	}
	defer file.Close()

	return t.parseRecords(file, filter)
}

func (t *TrustlessWorkConnector) parseRecords(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	var records []TrustlessWorkMilestoneRecord
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, fmt.Errorf("trustless-work connector error decoding json: %w", err)
	}

	var payments []types.InternalPayment
	for _, rec := range records {
		ts, err := time.Parse(time.RFC3339, rec.ReleasedAt)
		if err != nil {
			ts = time.Now()
		}

		if !filter.TimeStart.IsZero() && ts.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && ts.After(filter.TimeEnd) {
			continue
		}

		asset := rec.TokenSymbol
		if asset == "" {
			asset = "USDC"
		}

		payment := types.InternalPayment{
			ID:          fmt.Sprintf("TW-%s-M%d", rec.EscrowID, rec.MilestoneIndex),
			SourceApp:   "trustless-work",
			ReferenceID: rec.TransactionHash,
			Sender:      rec.ClientWallet,
			Recipient:   rec.ContractorWallet,
			Amount:      fmt.Sprintf("%.7f", rec.Amount),
			Asset:       asset,
			Timestamp:   ts,
			Status:      rec.State,
			Metadata: map[string]string{
				"escrow_id":         rec.EscrowID,
				"milestone_index":   fmt.Sprintf("%d", rec.MilestoneIndex),
				"contract_address":  rec.ContractAddress,
			},
		}

		payments = append(payments, payment)
	}

	return payments, nil
}
