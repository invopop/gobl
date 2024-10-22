package gb_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/gb"
	"github.com/stretchr/testify/assert"
)

func TestUKIdentifiers(t *testing.T) {
	tests := []struct {
		name          string
		idKey         cbc.Key
		initialCode   string
		expectedCode  string
		expectedError string
	}{
		{
			name:          "Normalize UTR - spaces removed",
			idKey:         gb.IdentityUTR,
			initialCode:   "  1234567890  ",
			expectedCode:  "1234567890",
			expectedError: "",
		},
		{
			name:          "Validate valid UTR",
			idKey:         gb.IdentityUTR,
			initialCode:   "1234567890",
			expectedCode:  "1234567890",
			expectedError: "",
		},
		{
			name:          "Validate invalid UTR - starts with 0",
			idKey:         gb.IdentityUTR,
			initialCode:   "0234567890",
			expectedCode:  "0234567890",
			expectedError: "code: invalid UTR format.",
		},
		{
			name:          "Normalize NINO - to uppercase",
			idKey:         gb.IdentityNINO,
			initialCode:   "ab123456c",
			expectedCode:  "AB123456C",
			expectedError: "",
		},
		{
			name:          "Validate valid NINO",
			idKey:         gb.IdentityNINO,
			initialCode:   "AB123456C",
			expectedCode:  "AB123456C",
			expectedError: "",
		},
		{
			name:          "Validate invalid NINO - disallowed prefix",
			idKey:         gb.IdentityNINO,
			initialCode:   "QQ123456Z",
			expectedCode:  "QQ123456Z",
			expectedError: "code: invalid NINO format.",
		},
		{
			name:          "Validate invalid NINO - incorrect format",
			idKey:         gb.IdentityNINO,
			initialCode:   "A123456C",
			expectedCode:  "A123456C",
			expectedError: "code: invalid NINO format.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{
				Key:  tt.idKey,
				Code: cbc.Code(tt.initialCode),
			}

			// Normalize the identifier
			gb.Normalize(id)

			// Check if the normalized code is as expected
			assert.Equal(t, tt.expectedCode, id.Code.String())

			// Validate the identifier
			err := gb.Validate(id)

			// Check if the error matches expected
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
