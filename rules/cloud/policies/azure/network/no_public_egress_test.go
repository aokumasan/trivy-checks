package network

import (
	"testing"

	defsecTypes "github.com/aquasecurity/defsec/pkg/types"

	"github.com/aquasecurity/defsec/pkg/state"

	"github.com/aquasecurity/defsec/pkg/providers/azure/network"
	"github.com/aquasecurity/defsec/pkg/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckNoPublicEgress(t *testing.T) {
	tests := []struct {
		name     string
		input    network.Network
		expected bool
	}{
		{
			name: "Security group outbound rule with wildcard destination address",
			input: network.Network{
				SecurityGroups: []network.SecurityGroup{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Rules: []network.SecurityGroupRule{
							{
								Metadata: defsecTypes.NewTestMetadata(),
								Allow:    defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
								Outbound: defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
								DestinationAddresses: []defsecTypes.StringValue{
									defsecTypes.String("*", defsecTypes.NewTestMetadata()),
								},
							},
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "Security group outbound rule with private destination address",
			input: network.Network{
				SecurityGroups: []network.SecurityGroup{
					{
						Metadata: defsecTypes.NewTestMetadata(),
						Rules: []network.SecurityGroupRule{
							{
								Metadata: defsecTypes.NewTestMetadata(),
								Allow:    defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
								Outbound: defsecTypes.Bool(true, defsecTypes.NewTestMetadata()),
								DestinationAddresses: []defsecTypes.StringValue{
									defsecTypes.String("10.0.0.0/16", defsecTypes.NewTestMetadata()),
								},
							},
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
			testState.Azure.Network = test.input
			results := CheckNoPublicEgress.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == CheckNoPublicEgress.LongID() {
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
