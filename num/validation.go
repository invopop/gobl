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

const (
	greaterThan = iota
	greaterEqualThan
	lessThan
	lessEqualThan
)

// Min checks if the value is greater than or equal to the provided amount or percentage
func Min(min interface{}) ThresholdRule {
	return ThresholdRule{
		threshold: interfaceToAmount(min),
		operator:  greaterEqualThan,
		err:       validation.ErrMinGreaterEqualThanRequired,
	}
}

// Max checks if the value is less than or equal to the provided amount or percentage
func Max(max interface{}) ThresholdRule {
	return ThresholdRule{
		threshold: interfaceToAmount(max),
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
		return a.Amount
	default:
		return Amount{}
	}
}

// Exclusive sets the comparison to exclude the boundary value.
func (r ThresholdRule) Exclusive() ThresholdRule {
	if r.operator == greaterEqualThan {
		r.operator = greaterThan
		r.err = validation.ErrMinGreaterThanRequired
	} else if r.operator == lessEqualThan {
		r.operator = lessThan
		r.err = validation.ErrMaxLessThanRequired
	}
	return r
}

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
	default:
		return cmp == -1 || cmp == 0
	}
}
