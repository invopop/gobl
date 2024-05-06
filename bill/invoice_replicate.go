package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/uuid"
)

// Replicate modifies the invoice's fields to ensure that it can be used as part
// of a replication process. A replicated invoice is one that contains the same
// base details like supplier, customer, lines, etc., but with updated identifiers
// and dates.
func (inv *Invoice) Replicate() error {
	inv.UUID = uuid.Empty
	inv.Code = ""
	inv.IssueDate = cal.Today()
	inv.ValueDate = nil
	inv.OperationDate = nil
	return nil
}
