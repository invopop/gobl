package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/uuid"
)

// Replicate modifies the payment's fields to ensure that it can be used as part
// of a replication process. A replicated payment is one that contains the same
// base details like supplier, customer, lines, etc., but with updated identifiers
// and dates.
func (pmt *Payment) Replicate() error {
	pmt.UUID = uuid.Empty
	pmt.Code = ""
	pmt.IssueDate = cal.Date{}
	pmt.IssueTime = nil
	pmt.ValueDate = nil
	return nil
}
