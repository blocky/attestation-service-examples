package price

import (
	"fmt"
	"sort"
	"time"
)

type Price struct {
	Value     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func TWAP(atTime time.Time, samples []Price) (float64, error) {
	if len(samples) == 0 {
		return 0, fmt.Errorf("no samples provided")
	}

	// Sort samples from latest to earliest
	sort.Slice(
		samples, func(i, j int) bool {
			return samples[i].Timestamp.After(samples[j].Timestamp)
		},
	)

	// Check that atTime is after the latest sample
	if atTime.Before(samples[0].Timestamp) {
		return 0, fmt.Errorf("atTime is before the latest sample")
	}

	var weightedSum, totalWeight float64

	for _, sample := range samples {
		timeDiff := atTime.Sub(sample.Timestamp).Seconds()
		weight := timeDiff
		weightedSum += sample.Value * weight
		totalWeight += weight
		atTime = sample.Timestamp
	}

	if totalWeight == 0 {
		return 0, fmt.Errorf("total weight is zero, cannot compute TWAP")
	}

	return weightedSum / totalWeight, nil
}
