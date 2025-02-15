package price

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeWeightedAverage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		samples  []Price
		expected float64
	}{
		{
			name: "multiple samples",
			samples: []Price{
				{Value: 100, Timestamp: now.Add(-1 * time.Hour)},
				{Value: 200, Timestamp: now.Add(-2 * time.Hour)},
				{Value: 300, Timestamp: now.Add(-3 * time.Hour)},
				{Value: 0, Timestamp: now.Add(-4 * time.Hour)},
			},
			expected: 200,
		},
		{
			name: "single sample",
			samples: []Price{
				{Value: 100, Timestamp: now.Add(-1 * time.Hour)},
			},
			expected: 100,
		},
		{
			name: "out of order samples",
			samples: []Price{
				{Value: 200, Timestamp: now.Add(-2 * time.Hour)},
				{Value: 100, Timestamp: now.Add(-1 * time.Hour)},
				{Value: 0, Timestamp: now.Add(-3 * time.Hour)},
			},
			expected: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			result, err := TWAP(tt.samples)
			require.NoError(t, err)

			// then
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("no samples", func(t *testing.T) {
		// when
		_, err := TWAP([]Price{})

		// then
		assert.ErrorContains(t, err, "no samples provided")
	})

	t.Run("no time difference", func(t *testing.T) {
		// given
		invalidSamples := []Price{
			{Value: 100, Timestamp: now},
			{Value: 200, Timestamp: now},
		}

		// when
		_, err := TWAP(invalidSamples)

		// then
		assert.ErrorContains(t, err, "total weight is zero, cannot compute TWAP")
	})
}
