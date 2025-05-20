package num

import (
	"github.com/invopop/validation"
)

// ThresholdRule is a validator for Amounts and Percentages
type ThresholdRule struct {
	threshold Amount
	operator  int
	err       validation.Error
}

var (
	// ErrIsZero indicates that the value is zero when it should not be.
	ErrIsZero = validation.NewError("validation_is_zero", "must not be zero")
)

const (
	greaterThan = iota
	greaterEqualThan
	lessThan
	lessEqualThan
	notZero
)

var (
	// Positive validates the that value is greater than or equal to zero.
	Positive = Min(MakeAmount(0, 0)).Exclusive()
	// Negative validates the value is less than or equal to zero.
	Negative = Max(MakeAmount(0, 0)).Exclusive()
	// NotZero validates that the value is not zero.
	NotZero = ThresholdRule{
		threshold: Amount{0, 0},
		operator:  notZero,
		err:       ErrIsZero,
	}
)

// Min checks if the value is greater than or equal to the provided amount or percentage
func Min(value any) ThresholdRule {
	return ThresholdRule{
		threshold: interfaceToAmount(value),
		operator:  greaterEqualThan,
		err:       validation.ErrMinGreaterEqualThanRequired,
	}
}

// Max checks if the value is less than or equal to the provided amount or percentage
func Max(value any) ThresholdRule {
	return ThresholdRule{
		threshold: interfaceToAmount(value),
		operator:  lessEqualThan,
		err:       validation.ErrMaxLessEqualThanRequired,
	}
}

func interfaceToAmount(val interface{}) Amount {
	val, isNil := validation.Indirect(val)
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
		r.err = validation.ErrMinGreaterThanRequired
	case lessEqualThan:
		r.operator = lessThan
		r.err = validation.ErrMaxLessThanRequired
	}
	return r
}

// Validate checks if the provided value confirms with the threshold
// rule.
func (r ThresholdRule) Validate(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	a := interfaceToAmount(value)
	if !r.compare(a) {
		return r.err.SetParams(map[string]interface{}{"threshold": r.threshold.String()})
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
	default:
		return cmp == -1 || cmp == 1
	}
}
