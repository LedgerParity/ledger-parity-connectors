package database

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

// ColumnMapping defines how database table columns map to InternalPayment fields.
type ColumnMapping struct {
	IDColumn        string `json:"id_column"`
	SenderColumn    string `json:"sender_column"`
	RecipientColumn string `json:"recipient_column"`
	AmountColumn    string `json:"amount_column"`
	CurrencyColumn  string `json:"currency_column"`
	TimestampColumn string `json:"timestamp_column"`
	StatusColumn    string `json:"status_column"`
	TxHashColumn    string `json:"tx_hash_column"`
}

// DefaultColumnMapping returns standard column names for payment tables.
func DefaultColumnMapping() ColumnMapping {
	return ColumnMapping{
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

	amountVal := formatAmountVal(row[d.mapping.AmountColumn])
	tsVal := parseTimestampVal(row[d.mapping.TimestampColumn])

	payment := types.InternalPayment{
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

func formatAmountVal(val interface{}) string {
	if val == nil {
		return "0.0000000"
	}
	switch v := val.(type) {
	case float64:
		return fmt.Sprintf("%.7f", v)
	case float32:
		return fmt.Sprintf("%.7f", float64(v))
	case int:
		return fmt.Sprintf("%.7f", float64(v))
	case int64:
		return fmt.Sprintf("%.7f", float64(v))
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return fmt.Sprintf("%.7f", f)
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func parseTimestampVal(val interface{}) time.Time {
	if val == nil {
		return time.Now()
	}
	switch v := val.(type) {
	case time.Time:
		return v
	case string:
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
		if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return t
		}
	case int64:
		return time.Unix(v, 0)
	}
	return time.Now()
}
