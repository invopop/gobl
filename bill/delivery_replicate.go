package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/uuid"
)

// Replicate modifies the delivery's fields to ensure that it can be used as part
// of a replication process. A replicated delivery is one that contains the same
// base details like supplier, customer, lines, etc., but with updated identifiers
// and dates.
func (dlv *Delivery) Replicate() error {
	dlv.UUID = uuid.Empty
	dlv.Code = ""
	dlv.IssueDate = cal.Date{}
	dlv.IssueTime = nil
	dlv.ValueDate = nil
	return nil
}
