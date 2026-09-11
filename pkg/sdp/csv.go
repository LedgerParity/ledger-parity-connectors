// Package sdp imports the payment CSV contract inspected at SDP 7.0.0.
// It does not contact SDP, submit payments, or establish export completeness.
package sdp

import (
	"encoding/csv"
	"fmt"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"io"
	"strings"
	"time"
)

const Revision = "14704274467d267b6c6679052d7251ef6050d118"
const Release = "7.0.0"
const Header = "ID,Amount,StellarTransactionID,Status,Type,Disbursement.ID,Asset.Code,Asset.Issuer,Wallet.Name,Receiver.ID,Receiver.PhoneNumber,Receiver.Email,Receiver.ExternalID,ReceiverWallet.Address,ReceiverWallet.Status,CreatedAt,UpdatedAt,ExternalPaymentID,CircleTransactionID,CircleTransactionType"

// Scope is an operator assertion obtained independently from settlement data.
// A single export must have a single known sending account. The interval bounds
// expected settlement, not CSV creation/update time. No completeness is inferred.
type Scope struct {
	Release         string    `json:"release"`
	Sender          string    `json:"sender"`
	SettlementStart time.Time `json:"settlement_start"`
	SettlementEnd   time.Time `json:"settlement_end"`
	Assertion       string    `json:"assertion"`
}

func (s Scope) Validate() error {
	if s.Release != Release || strings.TrimSpace(s.Sender) == "" || strings.TrimSpace(s.Assertion) == "" || s.SettlementStart.IsZero() || s.SettlementEnd.IsZero() || s.SettlementEnd.Before(s.SettlementStart) {
		return fmt.Errorf("SDP requires release %s, independently known sender, settlement interval and scope assertion", Release)
	}
	return nil
}

func Parse(r io.Reader, network, source string, scope Scope) ([]types.InternalPayment, error) {
	if err := scope.Validate(); err != nil {
		return nil, err
	}
	if network == "" || source == "" {
		return nil, fmt.Errorf("SDP network and deployment source required")
	}
	cr := csv.NewReader(r)
	headers, err := cr.Read()
	if err != nil {
		return nil, err
	}
	expected := strings.Split(Header, ",")
	if len(headers) != len(expected) {
		return nil, fmt.Errorf("SDP %s export header mismatch", Release)
	}
	cols := map[string]int{}
	allowed := map[string]bool{}
	for _, h := range expected {
		allowed[h] = true
	}
	for i, h := range headers {
		if _, ok := cols[h]; ok || !allowed[h] {
			return nil, fmt.Errorf("unexpected or duplicate SDP column %q", h)
		}
		cols[h] = i
	}
	result := []types.InternalPayment{}
	for n := 2; ; n++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("SDP row %d: invalid CSV", n)
		}
		v := func(k string) string { return row[cols[k]] }
		if v("CircleTransactionID") != "" || v("CircleTransactionType") != "" {
			return nil, fmt.Errorf("SDP row %d: Circle routes unsupported", n)
		}
		switch v("Type") {
		case "DISBURSEMENT", "DIRECT":
		default:
			return nil, fmt.Errorf("SDP row %d: unsupported payment type", n)
		}
		switch v("Status") {
		case "SUCCESS", "FAILED", "CANCELED", "DRAFT", "READY", "PENDING", "PAUSED":
		default:
			return nil, fmt.Errorf("SDP row %d: unknown status", n)
		}
		// Timestamps are checked but never used as settlement times or retained with PII.
		for _, k := range []string{"CreatedAt", "UpdatedAt"} {
			if _, err := time.Parse(time.RFC3339Nano, v(k)); err != nil {
				return nil, fmt.Errorf("SDP row %d: invalid %s", n, k)
			}
		}
		code, issuer := v("Asset.Code"), v("Asset.Issuer")
		kind := "credit_alphanum4"
		if issuer == "" && (code == "XLM" || code == "NATIVE") {
			kind = "native"
			code = "XLM"
		} else if len(code) > 4 {
			kind = "credit_alphanum12"
		}
		p := types.InternalPayment{ID: v("ID"), SourceApp: source, Network: network, OperationType: "payment", Sender: scope.Sender, Recipient: v("ReceiverWallet.Address"), Amount: v("Amount"), Asset: code, AssetType: kind, AssetIssuer: issuer, Status: strings.ToLower(v("Status")), ReferenceID: v("StellarTransactionID"), BusinessReference: v("ExternalPaymentID"), SettlementStart: scope.SettlementStart, SettlementEnd: scope.SettlementEnd}
		if p.Status == "success" && p.ReferenceID == "" {
			return nil, fmt.Errorf("SDP row %d: successful payment missing transaction reference", n)
		}
		if err := types.ValidateInternal(p); err != nil {
			return nil, fmt.Errorf("SDP row %d: %w", n, err)
		}
		// Preserve duplicate claims for the reconciler; never silently deduplicate.
		result = append(result, p)
	}
	return result, nil
}
