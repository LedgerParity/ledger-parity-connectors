package file

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/LedgerParity/ledger-parity-connectors/pkg/connector"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"io"
	"os"
	"strings"
	"time"
)

type FileConnector struct{ FilePath, Format, TargetName string }

func NewFileConnector(path, format, name string) *FileConnector {
	if format == "" {
		format = "json"
		if strings.HasSuffix(strings.ToLower(path), ".csv") {
			format = "csv"
		}
	}
	if name == "" {
		name = "file_export"
	}
	return &FileConnector{path, format, name}
}
func (f *FileConnector) Name() string { return f.TargetName }
func (f *FileConnector) FetchInternalPayments(ctx context.Context, filter connector.Filter) ([]types.InternalPayment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	input, err := os.Open(f.FilePath)
	if err != nil {
		return nil, err
	}
	defer input.Close()
	return f.Parse(ctx, input, filter)
}

// Parse consumes the same bytes a caller can hash for an evidence bundle.
func (f *FileConnector) Parse(ctx context.Context, input io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	var rows []types.InternalPayment
	var err error
	switch f.Format {
	case "json":
		rows, err = f.parseJSON(input, filter)
	case "csv":
		rows, err = f.parseCSV(input, filter)
	default:
		return nil, fmt.Errorf("unsupported file format %q", f.Format)
	}
	if err != nil {
		return nil, err
	}
	if err = ctx.Err(); err != nil {
		return nil, err
	}
	return rows, nil
}
func (f *FileConnector) parseJSON(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	var records []types.InternalPayment
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&records); err != nil {
		return nil, fmt.Errorf("invalid payment JSON: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("expected one JSON array")
	}
	if records == nil {
		return nil, fmt.Errorf("expected payment array, not null")
	}
	return f.validateFilter(records, filter)
}
func (f *FileConnector) validateFilter(records []types.InternalPayment, filter connector.Filter) ([]types.InternalPayment, error) {
	if filter.Limit < 0 || (!filter.TimeEnd.IsZero() && filter.TimeEnd.Before(filter.TimeStart)) {
		return nil, fmt.Errorf("invalid filter")
	}
	rows := []types.InternalPayment{}
	for i, p := range records {
		if p.SourceApp == "" {
			p.SourceApp = f.TargetName
		}
		if err := types.ValidateInternal(p); err != nil {
			return nil, fmt.Errorf("record %d: %w", i+1, err)
		}
		if !p.InWindow(filter.TimeStart, filter.TimeEnd) {
			if !p.SettlementStart.IsZero() && (filter.TimeStart.IsZero() || !p.SettlementEnd.Before(filter.TimeStart)) && (filter.TimeEnd.IsZero() || !p.SettlementStart.After(filter.TimeEnd)) {
				return nil, fmt.Errorf("record %d: settlement interval crosses filter boundary", i+1)
			}
			continue
		}
		if filter.Status != "" && !strings.EqualFold(p.Status, filter.Status) {
			continue
		}
		rows = append(rows, p)
	}
	if filter.Limit > 0 && len(rows) > filter.Limit {
		return nil, fmt.Errorf("filter limit would truncate export; use a narrower window or no limit")
	}
	return rows, nil
}
func (f *FileConnector) parseCSV(r io.Reader, filter connector.Filter) ([]types.InternalPayment, error) {
	reader := csv.NewReader(r)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	cols := map[string]int{}
	allowed := map[string]bool{}
	for _, h := range strings.Split("id,source_app,network,operation_type,operation_id,reference_id,business_reference,sender,recipient,amount,asset,asset_type,asset_issuer,asset_contract,timestamp,settlement_start,settlement_end,status", ",") {
		allowed[h] = true
	}
	for i, h := range headers {
		if !allowed[h] {
			return nil, fmt.Errorf("unknown CSV column %q", h)
		}
		if _, ok := cols[h]; ok {
			return nil, fmt.Errorf("duplicate CSV column %q", h)
		}
		cols[h] = i
	}
	rows := []types.InternalPayment{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		val := func(k string) string {
			if i, ok := cols[k]; ok {
				return row[i]
			}
			return ""
		}
		times := map[string]time.Time{}
		for _, key := range []string{"timestamp", "settlement_start", "settlement_end"} {
			if val(key) != "" {
				ts, err := time.Parse(time.RFC3339, val(key))
				if err != nil {
					return nil, fmt.Errorf("row %d: invalid %s", len(rows)+2, key)
				}
				times[key] = ts
			}
		}
		rows = append(rows, types.InternalPayment{ID: val("id"), SourceApp: val("source_app"), Network: val("network"), OperationType: val("operation_type"), OperationID: val("operation_id"), ReferenceID: val("reference_id"), BusinessReference: val("business_reference"), Sender: val("sender"), Recipient: val("recipient"), Amount: val("amount"), Asset: val("asset"), AssetType: val("asset_type"), AssetIssuer: val("asset_issuer"), AssetContract: val("asset_contract"), Timestamp: times["timestamp"], SettlementStart: times["settlement_start"], SettlementEnd: times["settlement_end"], Status: val("status")})
	}
	return f.validateFilter(rows, filter)
}
