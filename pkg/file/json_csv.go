package file

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

// FileConnector reads internal payment records from JSON or CSV files.
type FileConnector struct {
	FilePath   string
	Format     string // "json" or "csv"
	TargetName string
}

func NewFileConnector(filePath, format, targetName string) *FileConnector {
	if format == "" {
		if strings.HasSuffix(filePath, ".csv") {
			format = "csv"
		} else {
			format = "json"
		}
	}
	if targetName == "" {
		targetName = "file_export"
	}
	return &FileConnector{
		FilePath:   filePath,
		Format:     format,
		TargetName: targetName,
	}
}

func (f *FileConnector) Name() string {
	return f.TargetName
}

func (f *FileConnector) FetchInternalPayments(ctx context.Context, filter connector.Filter) ([]types.InternalPayment, error) {
	file, err := os.Open(f.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", f.FilePath, err)
	}
	defer file.Close()

	if strings.EqualFold(f.Format, "csv") {
		return f.parseCSV(file, filter)
	}
	return f.parseJSON(file, filter)
}

func (f *FileConnector) parseJSON(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	var records []types.InternalPayment
	if err := json.NewDecoder(r).Decode(&records); err != nil {
		return nil, fmt.Errorf("failed decoding JSON payments: %w", err)
	}

	var filtered []types.InternalPayment
	for _, p := range records {
		if p.SourceApp == "" {
			p.SourceApp = f.TargetName
		}
		if !filter.TimeStart.IsZero() && p.Timestamp.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && p.Timestamp.After(filter.TimeEnd) {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(p.Status, filter.Status) {
			continue
		}
		filtered = append(filtered, p)
	}

	if filter.Limit > 0 && len(filtered) > filter.Limit {
		filtered = filtered[:filter.Limit]
	}

	return filtered, nil
}

func (f *FileConnector) parseCSV(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed reading CSV header: %w", err)
	}

	colMap := make(map[string]int)
	for idx, h := range headers {
		colMap[strings.ToLower(strings.TrimSpace(h))] = idx
	}

	var payments []types.InternalPayment
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		getVal := func(keys ...string) string {
			for _, k := range keys {
				if idx, ok := colMap[k]; ok && idx < len(row) {
					return strings.TrimSpace(row[idx])
				}
			}
			return ""
		}

		tsStr := getVal("timestamp", "created_at", "date", "time")
		ts, err := time.Parse(time.RFC3339, tsStr)
		if err != nil {
			ts, err = time.Parse("2006-01-02 15:04:05", tsStr)
			if err != nil {
				ts = time.Now()
			}
		}

		if !filter.TimeStart.IsZero() && ts.Before(filter.TimeStart) {
			continue
		}
		if !filter.TimeEnd.IsZero() && ts.After(filter.TimeEnd) {
			continue
		}

		status := getVal("status", "state")
		if filter.Status != "" && !strings.EqualFold(status, filter.Status) {
			continue
		}

		asset := getVal("asset", "token", "currency")
		if asset == "" {
			asset = "XLM"
		}

		p := types.InternalPayment{
			ID:          getVal("id", "payment_id", "tx_id"),
			SourceApp:   f.TargetName,
			ReferenceID: getVal("reference_id", "ref_id", "hash", "transaction_hash"),
			Sender:      getVal("sender", "from", "source"),
			Recipient:   getVal("recipient", "to", "destination", "recipient_address"),
			Amount:      getVal("amount", "amt", "value"),
			Asset:       asset,
			Timestamp:   ts,
			Status:      status,
		}

		payments = append(payments, p)
	}

	if filter.Limit > 0 && len(payments) > filter.Limit {
		payments = payments[:filter.Limit]
	}

	return payments, nil
}
