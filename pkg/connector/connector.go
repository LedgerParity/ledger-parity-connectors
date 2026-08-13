package connector

import (
	"context"
	"time"

	"github.com/LedgerParity/ledger-parity-core/pkg/types"
)

// Filter specifies query criteria for fetching internal payment records.
type Filter struct {
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Limit     int       `json:"limit"`
	Status    string    `json:"status"`
}

// Connector defines the standard interface for target application adapters.
type Connector interface {
	Name() string
	FetchInternalPayments(ctx context.Context, filter Filter) ([]types.InternalPayment, error)
}
