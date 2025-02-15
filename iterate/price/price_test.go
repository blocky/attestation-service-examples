package price

import (
	"testing"
	"time"
)

func TestTimeWeightedAverage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		now      time.Time
		samples  []Price
		expected float64
	}{
		{
			name: "multiple samples",
			now:  now,
			samples: []Price{
				{Value: 100, Timestamp: now.Add(-1 * time.Hour)},
				{Value: 200, Timestamp: now.Add(-2 * time.Hour)},
				{Value: 300, Timestamp: now.Add(-3 * time.Hour)},
			},
			expected: 200,
		},
		{
			name: "single sample",
			now:  now,
			samples: []Price{
				{Value: 100, Timestamp: now.Add(-1 * time.Hour)},
			},
			expected: 100,
		},
		{
			name: "old samples",
			now:  now,
			samples: []Price{
				{Value: 100, Timestamp: now.Add(-3 * time.Hour)},
				{Value: 200, Timestamp: now.Add(-4 * time.Hour)},
			},
			expected: 125,
		},
		{
			name: "out of order samples",
			now:  now,
			samples: []Price{
				{Value: 200, Timestamp: now.Add(-4 * time.Hour)},
				{Value: 100, Timestamp: now.Add(-3 * time.Hour)},
			},
			expected: 125,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TWAP(tt.now, tt.samples)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}

	t.Run("no samples", func(t *testing.T) {
		_, err := TWAP(now, []Price{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "no samples provided" {
			t.Fatalf("expected 'no samples provided' error, got: %v", err)
		}
	})
}
