package trustlesswork

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

// TrustlessWorkMilestoneRecord is an experimental local example schema, not a verified upstream contract.
type TrustlessWorkMilestoneRecord struct {
	Network          string      `json:"network"`
	OperationType    string      `json:"operation_type"`
	OperationID      string      `json:"operation_id"`
	AssetType        string      `json:"asset_type"`
	AssetIssuer      string      `json:"asset_issuer"`
	AssetContract    string      `json:"asset_contract"`
	EscrowID         string      `json:"escrow_id"`
	MilestoneIndex   int         `json:"milestone_index"`
	ClientWallet     string      `json:"client_wallet"`
	ContractorWallet string      `json:"contractor_wallet"`
	Amount           json.Number `json:"amount"`
	TokenSymbol      string      `json:"token_symbol"` // e.g. "USDC", "XLM", "EURC"
	State            string      `json:"state"`        // "RELEASED", "DISPUTED", "RESOLVED", "REFUNDED"
	ReleasedAt       string      `json:"released_at"`
	TransactionHash  string      `json:"transaction_hash"`
	ContractAddress  string      `json:"contract_address"`
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
		amount, err := utils.ParsePaymentAmount(rec.Amount.String())
		if err != nil {
			return nil, err
		}
		ts, err := time.Parse(time.RFC3339, rec.ReleasedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
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
			Network: rec.Network, OperationType: rec.OperationType, OperationID: rec.OperationID, AssetType: rec.AssetType, AssetIssuer: rec.AssetIssuer, AssetContract: rec.AssetContract,
			ID:          fmt.Sprintf("TW-%s-M%d", rec.EscrowID, rec.MilestoneIndex),
			SourceApp:   "trustless-work",
			ReferenceID: rec.TransactionHash,
			Sender:      rec.ClientWallet,
			Recipient:   rec.ContractorWallet,
			Amount:      utils.FormatScaledAmount(amount, 7),
			Asset:       asset,
			Timestamp:   ts,
			Status:      rec.State,
			Metadata: map[string]string{
				"escrow_id":        rec.EscrowID,
				"milestone_index":  fmt.Sprintf("%d", rec.MilestoneIndex),
				"contract_address": rec.ContractAddress,
			},
		}

		if filter.Status != "" && !strings.EqualFold(payment.Status, filter.Status) {
			continue
		}
		payments = append(payments, payment)
	}

	if filter.Limit > 0 && len(payments) > filter.Limit {
		return nil, fmt.Errorf("limit would truncate input")
	}
	return payments, nil
}
