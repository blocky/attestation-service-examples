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

func TWAP(now time.Time, samples []Price) (float64, error) {
	if len(samples) == 0 {
		return 0, fmt.Errorf("no samples provided")
	}

	// Sort samples from latest to earliest
	sort.Slice(
		samples, func(i, j int) bool {
			return samples[i].Timestamp.After(samples[j].Timestamp)
		},
	)

	var weightedSum, totalWeight float64

	for _, sample := range samples {
		timeDiff := now.Sub(sample.Timestamp).Seconds()
		weight := timeDiff
		weightedSum += sample.Value * weight
		totalWeight += weight
		now = sample.Timestamp
	}

	return weightedSum / totalWeight, nil
}
