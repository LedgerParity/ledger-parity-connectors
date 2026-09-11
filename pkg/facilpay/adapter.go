package facilpay

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

// FacilPayTransactionRecord is an experimental local example schema, not a verified upstream contract.
type FacilPayTransactionRecord struct {
	Network           string      `json:"network"`
	OperationType     string      `json:"operation_type"`
	OperationID       string      `json:"operation_id"`
	AssetType         string      `json:"asset_type"`
	AssetIssuer       string      `json:"asset_issuer"`
	AssetContract     string      `json:"asset_contract"`
	TxUUID            string      `json:"tx_uuid"`
	MerchantKey       string      `json:"merchant_key"`
	CustomerPublicKey string      `json:"customer_public_key"`
	Amount            json.Number `json:"amount"`
	Token             string      `json:"token"` // e.g. "USDC", "XLM"
	State             string      `json:"state"` // "SETTLED", "SETTLEMENT_PENDING", "FAILED"
	CreatedAt         string      `json:"created_at"`
	StellarMemo       string      `json:"stellar_memo"`
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
		amount, err := utils.ParsePaymentAmount(rec.Amount.String())
		if err != nil {
			return nil, err
		}
		ts, err := time.Parse(time.RFC3339, rec.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
		}

		if !filter.TimeStart.IsZero() && ts.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && ts.After(filter.TimeEnd) {
			continue
		}

		p := types.InternalPayment{
			Network: rec.Network, OperationType: rec.OperationType, OperationID: rec.OperationID, AssetType: rec.AssetType, AssetIssuer: rec.AssetIssuer, AssetContract: rec.AssetContract,
			ID:          rec.TxUUID,
			SourceApp:   "facilpay",
			ReferenceID: rec.StellarMemo,
			Sender:      rec.CustomerPublicKey,
			Recipient:   rec.MerchantKey,
			Amount:      utils.FormatScaledAmount(amount, 7),
			Asset:       rec.Token,
			Timestamp:   ts,
			Status:      rec.State,
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
