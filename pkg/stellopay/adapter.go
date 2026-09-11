package stellopay

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"github.com/LedgerParity/ledger-parity-core/pkg/utils"
)

// StellopayPaymentRecord is an experimental local example schema, not a verified upstream contract.
type StellopayPaymentRecord struct {
	Network        string      `json:"network"`
	OperationType  string      `json:"operation_type"`
	OperationID    string      `json:"operation_id"`
	AssetType      string      `json:"asset_type"`
	AssetIssuer    string      `json:"asset_issuer"`
	AssetContract  string      `json:"asset_contract"`
	PaymentID      string      `json:"payment_id"`
	BatchID        string      `json:"batch_id"`
	EmployeeWallet string      `json:"employee_wallet"`
	EmployerWallet string      `json:"employer_wallet"`
	AmountXLM      json.Number `json:"amount_xlm"`
	Currency       string      `json:"currency"`
	Status         string      `json:"status"` // "COMPLETED", "SUBMITTED", "FAILED"
	ProcessedAt    string      `json:"processed_at"`
	StellarTxHash  string      `json:"stellar_tx_hash"`
}

// StellopayConnector adapts Stellopay's payment records into normalized InternalPayment structs.
type StellopayConnector struct {
	SourcePath string
}

func NewStellopayConnector(sourcePath string) *StellopayConnector {
	return &StellopayConnector{SourcePath: sourcePath}
}

func (s *StellopayConnector) Name() string {
	return "stellopay"
}

func (s *StellopayConnector) FetchInternalPayments(ctx context.Context, filter connector.Filter) ([]types.InternalPayment, error) {
	file, err := os.Open(s.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("stellopay connector error opening %s: %w", s.SourcePath, err)
	}
	defer file.Close()

	return s.parseRecords(file, filter)
}

func (s *StellopayConnector) parseRecords(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	var records []StellopayPaymentRecord
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, fmt.Errorf("stellopay connector error decoding json: %w", err)
	}

	var payments []types.InternalPayment
	for _, rec := range records {
		amount, err := utils.ParsePaymentAmount(rec.AmountXLM.String())
		if err != nil {
			return nil, err
		}
		ts, err := time.Parse(time.RFC3339, rec.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
		}

		if !filter.TimeStart.IsZero() && ts.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && ts.After(filter.TimeEnd) {
			continue
		}

		asset := rec.Currency
		if asset == "" {
			asset = "XLM"
		}

		p := types.InternalPayment{
			Network: rec.Network, OperationType: rec.OperationType, OperationID: rec.OperationID, AssetType: rec.AssetType, AssetIssuer: rec.AssetIssuer, AssetContract: rec.AssetContract,
			ID:          rec.PaymentID,
			SourceApp:   "stellopay",
			ReferenceID: rec.StellarTxHash,
			Sender:      rec.EmployerWallet,
			Recipient:   rec.EmployeeWallet,
			Amount:      utils.FormatScaledAmount(amount, 7),
			Asset:       asset,
			Timestamp:   ts,
			Status:      rec.Status,
			Metadata: map[string]string{
				"batch_id": rec.BatchID,
			},
		}

		if filter.Status != "" && !strings.EqualFold(p.Status, filter.Status) {
			continue
		}
		payments = append(payments, p)
	}

	if filter.Limit > 0 && len(payments) > filter.Limit {
		return nil, fmt.Errorf("limit would truncate input")
	}
	return payments, nil
}
