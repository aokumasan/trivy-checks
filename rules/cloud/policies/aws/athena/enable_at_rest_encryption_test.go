package athena

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/aws/athena"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckEnableAtRestEncryption(t *testing.T) {
	tests := []struct {
		name     string
		input    athena.Athena
		expected bool
	}{
		{
			name: "AWS Athena database unencrypted",
			input: athena.Athena{
				Databases: []athena.Database{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Encryption: athena.EncryptionConfiguration{
							Metadata: defsecTypes.NewTestMetadata(),
							Type:     defsecTypes.String(athena.EncryptionTypeNone, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "AWS Athena workgroup unencrypted",
			input: athena.Athena{
				Workgroups: []athena.Workgroup{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Encryption: athena.EncryptionConfiguration{
							Metadata: defsecTypes.NewTestMetadata(),
							Type:     defsecTypes.String(athena.EncryptionTypeNone, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "AWS Athena database and workgroup encrypted",
			input: athena.Athena{
				Databases: []athena.Database{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Encryption: athena.EncryptionConfiguration{
							Metadata: defsecTypes.NewTestMetadata(),
							Type:     defsecTypes.String(athena.EncryptionTypeSSEKMS, defsecTypes.NewTestMetadata()),
						},
					},
				},
				Workgroups: []athena.Workgroup{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Encryption: athena.EncryptionConfiguration{
							Metadata: defsecTypes.NewTestMetadata(),
							Type:     defsecTypes.String(athena.EncryptionTypeSSEKMS, defsecTypes.NewTestMetadata()),
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.Athena = test.input
			results := CheckEnableAtRestEncryption.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckEnableAtRestEncryption.LongID() {
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
