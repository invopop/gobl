package num

import (
	"fmt"

	"github.com/invopop/gobl/rules"
)

// ThresholdRule is a validator for Amounts and Percentages
type ThresholdRule struct {
	desc      string
	threshold Amount
	operator  int
	err       string // format string with %s for threshold
}

const (
	greaterThan = iota
	greaterEqualThan
	lessThan
	lessEqualThan
	notZero
	equals
)

var (
	// Positive validates the that value is greater than zero.
	Positive = Min(MakeAmount(0, 0)).Exclusive()
	// ZeroOrPositive validates the that value is greater than or equal to zero.
	ZeroOrPositive = Min(MakeAmount(0, 0))
	// Negative validates the value is less than zero.
	Negative = Max(MakeAmount(0, 0)).Exclusive()
	// ZeroOrNegative validates the value is less than or equal to zero.
	ZeroOrNegative = Max(MakeAmount(0, 0))
	// NotZero validates that the value is not zero.
	NotZero = ThresholdRule{
		desc:      "not zero",
		threshold: Amount{0, 0},
		operator:  notZero,
		err:       "must not be zero",
	}
)

// Check returns true if the value satisfies the threshold rule.
func (r ThresholdRule) Check(value any) bool {
	return r.Validate(value) == nil
}

// String returns the string representation of the threshold rule,
// which is part of the rules.Test interface.
func (r ThresholdRule) String() string {
	return r.desc
}

// Min checks if the value is greater than or equal to the provided amount or percentage
func Min(value any) ThresholdRule {
	return ThresholdRule{
		desc:      fmt.Sprintf("min %s", value),
		threshold: interfaceToAmount(value),
		operator:  greaterEqualThan,
		err:       "must be no less than %s",
	}
}

// Max checks if the value is less than or equal to the provided amount or percentage
func Max(value any) ThresholdRule {
	return ThresholdRule{
		desc:      fmt.Sprintf("max %s", value),
		threshold: interfaceToAmount(value),
		operator:  lessEqualThan,
		err:       "must be no greater than %s",
	}
}

// Equals checks if the value is equal to the provided amount or percentage
func Equals(value any) ThresholdRule {
	return ThresholdRule{
		desc:      fmt.Sprintf("equals %s", value),
		threshold: interfaceToAmount(value),
		operator:  equals,
		err:       "must be equal to %s",
	}
}

func interfaceToAmount(val interface{}) Amount {
	val, isNil := rules.Indirect(val)
	if isNil {
		return Amount{}
	}
	switch a := val.(type) {
	case Amount:
		return a
	case Percentage:
		return a.amount
	default:
		return Amount{}
	}
}

// Exclusive sets the comparison to exclude the boundary value.
func (r ThresholdRule) Exclusive() ThresholdRule {
	switch r.operator {
	case greaterEqualThan:
		r.operator = greaterThan
		r.err = "must be greater than %s"
	case lessEqualThan:
		r.operator = lessThan
		r.err = "must be less than %s"
	}
	return r
}

// Validate checks if the provided value confirms with the threshold
// rule.
func (r ThresholdRule) Validate(value interface{}) error {
	value, isNil := rules.Indirect(value)
	if isNil || rules.IsEmpty(value) {
		return nil
	}

	a := interfaceToAmount(value)
	if !r.compare(a) {
		return fmt.Errorf(r.err, r.threshold.String())
	}

	return nil
}

func (r ThresholdRule) compare(value Amount) bool {
	cmp := value.Compare(r.threshold)
	switch r.operator {
	case greaterThan:
		return cmp == 1
	case greaterEqualThan:
		return cmp == 1 || cmp == 0
	case lessThan:
		return cmp == -1
	case lessEqualThan:
		return cmp == -1 || cmp == 0
	case equals:
		return cmp == 0
	default:
		return cmp == -1 || cmp == 1
	}
}
