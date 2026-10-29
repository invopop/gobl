package mydata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncomeExtValidatorsWithNonExtensionValues(t *testing.T) {
	assert.True(t, incomeExtPairValid("not extensions"))
	assert.True(t, incomeExtOtherInfoValid("not extensions"))
}
