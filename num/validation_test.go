package num

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMin(t *testing.T) {
	tests := []struct {
		tag       string
		threshold interface{}
		exclusive bool
		value     interface{}
		err       string
	}{
		{"t1.1", MakeAmount(0, 0), false, MakeAmount(10, 0), ""},
		{"t1.2", MakeAmount(2, 0), false, MakeAmount(10, 0), ""},
		{"t1.3", MakeAmount(-2, 0), false, MakeAmount(10, 0), ""},
		{"t1.4", MakeAmount(10, 0), false, MakeAmount(10, 0), ""},
		{"t1.5", MakeAmount(9, 0), true, MakeAmount(10, 0), ""},
		{"t1.6", NewAmount(0, 0), true, MakeAmount(10, 0), ""},
		{"t1.7", NewAmount(0, 0), true, NewAmount(10, 0), ""},
		{"t1.8", MakeAmount(0, 0), true, NewAmount(10, 0), ""},
		{"t1.9", MakeAmount(10, 0), true, MakeAmount(10, 0), "must be greater than 10"},
		{"t1.10", MakeAmount(11, 0), false, MakeAmount(10, 0), "must be no less than 11"},
		{"t1.11", MakeAmount(-1, 0), false, MakeAmount(-2, 0), "must be no less than -1"},
		{"t2.1", MakePercentage(0, 0), false, MakePercentage(10, 0), ""},
		{"t2.2", MakePercentage(2, 0), false, MakePercentage(10, 0), ""},
		{"t2.3", MakePercentage(2, 0), false, MakePercentage(1, 0), "must be no less than 2"},
		{"t3.1", MakeAmount(0, 0), false, "", ""},   // ignore empty
		{"t3.2", nil, false, MakeAmount(10, 0), ""}, // same as zero
		{"t3.3", "", false, MakeAmount(10, 0), ""},  // same as zero
		{"t3.4", "", false, MakeAmount(-1, 0), "must be no less than 0"},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			r := Min(test.threshold)
			if test.exclusive {
				r = r.Exclusive()
			}
			err := r.Validate(test.value)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), test.err)
				}
			}
		})
	}
}

func TestValidationMax(t *testing.T) {
	tests := []struct {
		tag       string
		threshold interface{}
		exclusive bool
		value     interface{}
		err       string
	}{
		{"t1.1", MakeAmount(10, 0), false, MakeAmount(5, 0), ""},
		{"t1.2", MakeAmount(0, 0), false, MakeAmount(-2, 0), ""},
		{"t1.3", MakeAmount(-2, 0), false, MakeAmount(-5, 0), ""},
		{"t1.4", MakeAmount(10, 0), false, MakeAmount(10, 0), ""},
		{"t1.5", MakeAmount(11, 0), true, MakeAmount(10, 0), ""},
		{"t1.6", NewAmount(10, 0), true, MakeAmount(0, 0), ""},
		{"t1.9", MakeAmount(10, 0), true, MakeAmount(10, 0), "must be less than 10"},
		{"t1.10", MakeAmount(9, 0), false, MakeAmount(10, 0), "must be no greater than 9"},
		{"t1.10", MakeAmount(-5, 0), false, MakeAmount(-1, 0), "must be no greater than -5"},
		{"t2.1", MakePercentage(10, 0), false, MakePercentage(1, 0), ""},
		{"t2.2", MakePercentage(2, 0), false, MakePercentage(-1, 0), ""},
		{"t2.3", MakePercentage(1, 0), false, MakePercentage(2, 0), "must be no greater than 1"},
		{"t3.1", MakeAmount(0, 0), false, "", ""},   // ignore empty
		{"t3.2", nil, false, MakeAmount(-1, 0), ""}, // same as zero
		{"t3.3", "", false, MakeAmount(-1, 0), ""},  // same as zero
		{"t3.4", "", false, MakeAmount(10, 0), "must be no greater than 0"},
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			r := Max(test.threshold)
			if test.exclusive {
				r = r.Exclusive()
			}
			err := r.Validate(test.value)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), test.err)
				}
			}
		})
	}
}

func TestValidationFixedRules(t *testing.T) {
	tests := []struct {
		tag   string
		rule  ThresholdRule
		value interface{}
		err   string
	}{
		{"t1.1", Positive, MakeAmount(10, 0), ""},
		{"t1.2", Positive, MakeAmount(0, 0), "must be greater than 0"},
		{"t1.3", Positive, MakeAmount(-10, 0), "must be greater than 0"},
		{"t1.4", Positive, nil, ""}, // ignore empty
		{"t2.1", Negative, MakeAmount(-10, 0), ""},
		{"t2.2", Negative, MakeAmount(0, 0), "must be less than 0"},
		{"t2.3", Negative, MakeAmount(10, 0), "must be less than 0"},
		{"t2.4", Negative, nil, ""}, // ignore empty
		{"t3.1", NotZero, MakeAmount(10, 0), ""},
		{"t3.2", NotZero, MakeAmount(-10, 0), ""},
		{"t3.3", NotZero, MakeAmount(0, 0), "must not be zero"},
		{"t3.4", NotZero, nil, ""}, // out of scope!
	}

	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			err := test.rule.Validate(test.value)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), test.err)
				}
			}
		})
	}
}
