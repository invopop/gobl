package rules_test

import (
	"regexp"
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var codeFormat = regexp.MustCompile(`^[A-Z0-9]+(-[A-Z0-9]+)*$`)

func collectAssertionIDs(s *rules.Set) []rules.Code {
	var ids []rules.Code
	for _, a := range s.Assert {
		ids = append(ids, a.ID)
	}
	for _, sub := range s.Subsets {
		ids = append(ids, collectAssertionIDs(sub)...)
	}
	return ids
}

func TestFaultCodeIntegrity(t *testing.T) {
	var allIDs []rules.Code
	for _, s := range rules.AllSets() {
		allIDs = append(allIDs, collectAssertionIDs(s)...)
	}
	require.NotEmpty(t, allIDs, "expected at least one registered assertion")

	t.Run("format", func(t *testing.T) {
		for _, id := range allIDs {
			assert.Truef(t, codeFormat.MatchString(string(id)),
				"fault code %q does not match expected format [A-Z0-9]+(-[A-Z0-9]+)*", id)
		}
	})

	t.Run("unique", func(t *testing.T) {
		seen := make(map[rules.Code]int)
		for _, id := range allIDs {
			seen[id]++
		}
		for id, count := range seen {
			assert.Equalf(t, 1, count,
				"fault code %q is registered %d times, expected exactly once", id, count)
		}
	})
}
