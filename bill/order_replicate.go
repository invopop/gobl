package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/uuid"
)

// Replicate modifies the order's fields to ensure that it can be used as part
// of a replication process. A replicated order is one that contains the same
// base details like supplier, customer, lines, etc., but with updated identifiers
// and dates.
func (ord *Order) Replicate() error {
	ord.UUID = uuid.Empty
	ord.Code = ""
	ord.IssueDate = cal.Date{}
	ord.IssueTime = nil
	ord.ValueDate = nil
	ord.OperationDate = nil
	return nil
}
