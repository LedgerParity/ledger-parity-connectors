package sdp

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/LedgerParity/ledger-parity-core/pkg/engine"
	"github.com/LedgerParity/ledger-parity-core/pkg/types"
	"strings"
	"testing"
	"time"
)

func fixture(edit map[string]string, duplicate bool) string {
	h := strings.Split(Header, ",")
	v := map[string]string{"ID": "p1", "Amount": "922337203685.4775807", "StellarTransactionID": "tx", "Status": "SUCCESS", "Type": "DIRECT", "Asset.Code": "XLM", "ReceiverWallet.Address": "recipient", "CreatedAt": "2026-09-01T10:00:00Z", "UpdatedAt": "2026-09-01T13:00:00Z", "ExternalPaymentID": "invoice-1", "Receiver.Email": "private@example.invalid", "Receiver.PhoneNumber": "private-phone"}
	for k, x := range edit {
		v[k] = x
	}
	row := make([]string, len(h))
	for i, k := range h {
		row[i] = v[k]
	}
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.Write(h)
	w.Write(row)
	if duplicate {
		w.Write(row)
	}
	w.Flush()
	return b.String()
}
func scope() Scope {
	a, _ := time.Parse(time.RFC3339, "2026-09-01T11:00:00Z")
	return Scope{Release: Release, Sender: "sender", SettlementStart: a, SettlementEnd: a.Add(2 * time.Hour), Assertion: "synthetic one-account export; independently bounded interval"}
}
func TestContractAndMatching(t *testing.T) {
	s := scope()
	ps, err := Parse(strings.NewReader(fixture(nil, false)), "network", "deployment", s)
	if err != nil {
		t.Fatal(err)
	}
	p := ps[0]
	raw, _ := json.Marshal(ps)
	if bytes.Contains(raw, []byte("private")) || p.BusinessReference != "invoice-1" || p.ReferenceID != "tx" || !p.Timestamp.IsZero() {
		t.Fatal("provenance/privacy mapping", string(raw))
	}
	o := types.OnChainPayment{Network: p.Network, OperationType: "payment", AssetType: p.AssetType, AssetCode: p.Asset, Amount: p.Amount, Account: p.Sender, Destination: p.Recipient, OperationID: "1", TransactionHash: "tx", Timestamp: s.SettlementStart, Successful: true}
	opts := engine.ReconcileOptions{Coverage: types.Coverage{Network: p.Network, Accounts: []string{p.Sender}, Start: s.SettlementStart, End: s.SettlementEnd, Complete: true}}
	run := func(ps []types.InternalPayment, os []types.OnChainPayment) *types.DiscrepancyReport {
		return engine.NewReconciler(opts).Reconcile("sdp", s.SettlementStart, s.SettlementEnd, ps, os)
	}
	if r := run(ps, []types.OnChainPayment{o}); r.TotalMatched != 1 {
		t.Fatal(r)
	}
	o2 := o
	o2.OperationID = "2"
	if r := run(ps, []types.OnChainPayment{o, o2}); r.Results[0].Status != types.MatchUnknown {
		t.Fatal("batch ambiguity lost")
	}
	dup, err := Parse(strings.NewReader(fixture(nil, true)), "network", "deployment", s)
	if err != nil || len(dup) != 2 {
		t.Fatal(err)
	}
	if r := run(dup, []types.OnChainPayment{o}); r.DiscrepancyCounts[types.DiscrepancyDuplicateInternal] != 2 {
		t.Fatal("duplicate suppressed")
	}
	o.AssetType = "credit_alphanum4"
	o.AssetIssuer = "issuer"
	if r := run(ps, []types.OnChainPayment{o}); r.TotalMatched != 0 || r.DiscrepancyCounts[types.DiscrepancyOrphanedOnChain] != 0 {
		t.Fatal("asset alias or false orphan")
	}
}
func TestRejectUnsupportedAndMalformed(t *testing.T) {
	for _, edit := range []map[string]string{{"Amount": "1.00000001"}, {"ReceiverWallet.Address": ""}, {"CircleTransactionID": "circle"}, {"StellarTransactionID": ""}, {"Status": "NEW_STATUS"}, {"Type": "OTHER"}, {"CreatedAt": "yesterday"}, {"Asset.Code": "USDC", "Asset.Issuer": ""}} {
		if _, err := Parse(strings.NewReader(fixture(edit, false)), "network", "source", scope()); err == nil {
			t.Fatal("accepted", edit)
		}
	}
	s := scope()
	s.Sender = ""
	if _, err := Parse(strings.NewReader(fixture(nil, false)), "network", "source", s); err == nil {
		t.Fatal("missing scope")
	}
	raw := strings.Replace(fixture(nil, false), "Amount,", "ID,", 1)
	if _, err := Parse(strings.NewReader(raw), "network", "source", scope()); err == nil {
		t.Fatal("duplicate header")
	}
	ps, err := Parse(strings.NewReader(fixture(map[string]string{"Asset.Code": "NATIVE"}, false)), "network", "source", scope())
	if err != nil || ps[0].Asset != "XLM" {
		t.Fatal("native mapping", err)
	}
}
