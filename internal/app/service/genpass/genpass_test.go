package genpass

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenPass_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		expectedLen   int
		expectedError error
	}{
		{
			name: "normal case - length 12",

			expectedLen:   12,
			expectedError: nil,
		},
		{
			name: "normal case - length 16",

			expectedLen:   16,
			expectedError: nil,
		},
		{
			name: "normal case - length 1000",

			expectedLen:   10000,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			testSvc := NewGenPass()
			password, err := testSvc.GeneratePassword(tc.expectedLen)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedLen, len(password))
		})
	}

}
