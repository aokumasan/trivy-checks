package iam

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/iam"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckRequireUppercaseInPasswords(t *testing.T) {
	tests := []struct {
		name     string
		input    iam.IAM
		expected bool
	}{
		{
			name: "IAM password policy uppercase not required",
			input: iam.IAM{
				PasswordPolicy: iam.PasswordPolicy{
					Metadata:         defsecTypes.NewTestMetadata(),
					RequireUppercase: defsecTypes.Bool(false, defsecTypes.NewTestMetadata()),
				},
			},
			expected: true,
		},
		{
			name: "IAM password policy uppercase required",
			input: iam.IAM{
				PasswordPolicy: iam.PasswordPolicy{
					Metadata:         defsecTypes.NewTestMetadata(),
					RequireUppercase: defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.IAM = test.input
			results := CheckRequireUppercaseInPasswords.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckRequireUppercaseInPasswords.GetRule().LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
