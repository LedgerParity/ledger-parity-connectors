package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LedgerParity/ledger-parity-core/pkg/utils"
	"strconv"
	"time"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

// ColumnMapping defines how database table columns map to InternalPayment fields.
type ColumnMapping struct {
	NetworkColumn       string `json:"network_column"`
	OperationTypeColumn string `json:"operation_type_column"`
	OperationIDColumn   string `json:"operation_id_column"`
	AssetTypeColumn     string `json:"asset_type_column"`
	AssetIssuerColumn   string `json:"asset_issuer_column"`
	IDColumn            string `json:"id_column"`
	SenderColumn        string `json:"sender_column"`
	RecipientColumn     string `json:"recipient_column"`
	AmountColumn        string `json:"amount_column"`
	CurrencyColumn      string `json:"currency_column"`
	TimestampColumn     string `json:"timestamp_column"`
	StatusColumn        string `json:"status_column"`
	TxHashColumn        string `json:"tx_hash_column"`
}

// DefaultColumnMapping returns standard column names for payment tables.
func DefaultColumnMapping() ColumnMapping {
	return ColumnMapping{NetworkColumn: "network", OperationTypeColumn: "operation_type", OperationIDColumn: "operation_id", AssetTypeColumn: "asset_type", AssetIssuerColumn: "asset_issuer",
		IDColumn:        "id",
		SenderColumn:    "sender_address",
		RecipientColumn: "recipient_address",
		AmountColumn:    "amount",
		CurrencyColumn:  "currency",
		TimestampColumn: "created_at",
		StatusColumn:    "status",
		TxHashColumn:    "tx_hash",
	}
}

// RowData represents a generic key-value row returned from a database query.
type RowData map[string]interface{}

// DatabaseConnector converts relational database query rows into normalized InternalPayments.
type DatabaseConnector struct {
	name       string
	sourceApp  string
	mapping    ColumnMapping
	fetchQuery func(ctx context.Context, filter connector.Filter) ([]RowData, error)
}

// NewDatabaseConnector creates a new database connector with custom column mappings.
func NewDatabaseConnector(name, sourceApp string, mapping ColumnMapping, fetcher func(ctx context.Context, filter connector.Filter) ([]RowData, error)) *DatabaseConnector {
	if name == "" {
		name = "database-connector"
	}
	if sourceApp == "" {
		sourceApp = "database"
	}
	return &DatabaseConnector{
		name:       name,
		sourceApp:  sourceApp,
		mapping:    mapping,
		fetchQuery: fetcher,
	}
}

func (d *DatabaseConnector) Name() string {
	return d.name
}

func (d *DatabaseConnector) FetchInternalPayments(ctx context.Context, filter connector.Filter) ([]types.InternalPayment, error) {
	if d.fetchQuery == nil {
		return nil, fmt.Errorf("database connector %s has no query fetcher configured", d.name)
	}

	rows, err := d.fetchQuery(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed executing query for connector %s: %w", d.name, err)
	}

	var payments []types.InternalPayment
	for i, row := range rows {
		payment, err := d.MapRowToPayment(row)
		if err != nil {
			return nil, fmt.Errorf("error mapping row %d in connector %s: %w", i, d.name, err)
		}

		if !filter.TimeStart.IsZero() && payment.Timestamp.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && payment.Timestamp.After(filter.TimeEnd) {
			continue
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

// MapRowToPayment converts a single database query row into an InternalPayment.
func (d *DatabaseConnector) MapRowToPayment(row RowData) (types.InternalPayment, error) {
	idVal := getStringVal(row, d.mapping.IDColumn)
	if idVal == "" {
		return types.InternalPayment{}, fmt.Errorf("missing required ID column: %s", d.mapping.IDColumn)
	}

	senderVal := getStringVal(row, d.mapping.SenderColumn)
	recipientVal := getStringVal(row, d.mapping.RecipientColumn)
	statusVal := getStringVal(row, d.mapping.StatusColumn)
	txHashVal := getStringVal(row, d.mapping.TxHashColumn)
	currVal := getStringVal(row, d.mapping.CurrencyColumn)
	if currVal == "" {
		currVal = "XLM"
	}

	amountVal, err := formatAmountVal(row[d.mapping.AmountColumn])
	if err != nil {
		return types.InternalPayment{}, err
	}
	tsVal, err := parseTimestampVal(row[d.mapping.TimestampColumn])
	if err != nil {
		return types.InternalPayment{}, err
	}

	payment := types.InternalPayment{Network: getStringVal(row, d.mapping.NetworkColumn), OperationType: getStringVal(row, d.mapping.OperationTypeColumn), OperationID: getStringVal(row, d.mapping.OperationIDColumn), AssetType: getStringVal(row, d.mapping.AssetTypeColumn), AssetIssuer: getStringVal(row, d.mapping.AssetIssuerColumn),
		ID:          idVal,
		SourceApp:   d.sourceApp,
		ReferenceID: txHashVal,
		Sender:      senderVal,
		Recipient:   recipientVal,
		Amount:      amountVal,
		Asset:       currVal,
		Timestamp:   tsVal,
		Status:      statusVal,
		Metadata: map[string]string{
			"connector_type": "relational_db",
		},
	}

	return payment, nil
}

func getStringVal(row RowData, key string) string {
	if key == "" {
		return ""
	}
	if v, exists := row[key]; exists && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func formatAmountVal(val interface{}) (string, error) {
	var s string
	switch v := val.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	case json.Number:
		s = v.String()
	case int:
		s = strconv.Itoa(v)
	case int64:
		s = strconv.FormatInt(v, 10)
	default:
		return "", fmt.Errorf("amount requires decimal text or integer, never floating point")
	}
	n, err := utils.ParsePaymentAmount(s)
	if err != nil {
		return "", err
	}
	return utils.FormatScaledAmount(n, 7), nil
}
func parseTimestampVal(val interface{}) (time.Time, error) {
	switch v := val.(type) {
	case time.Time:
		if !v.IsZero() {
			return v, nil
		}
	case string:
		return time.Parse(time.RFC3339, v)
	case int64:
		return time.Unix(v, 0).UTC(), nil
	}
	return time.Time{}, fmt.Errorf("invalid or missing timestamp")
}
