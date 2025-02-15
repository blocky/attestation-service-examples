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

func TWAP(samples []Price) (float64, error) {
	if len(samples) == 0 {
		return 0, fmt.Errorf("no samples provided")
	}
	if len(samples) == 1 {
		return samples[0].Value, nil
	}

	// Sort samples from latest to earliest
	sort.Slice(
		samples, func(i, j int) bool {
			return samples[i].Timestamp.After(samples[j].Timestamp)
		},
	)

	var weightedSum, totalWeight float64

	// IMPORTANT: The value of the last sample is not included in the calculation
	// because it doesn't have a next sample to compare with. However, its
	// timestamp is used to calculate the weight of the previous sample.
	prev := samples[0]
	for _, next := range samples[1:] {
		timeDiff := prev.Timestamp.Sub(next.Timestamp).Microseconds()
		weight := float64(timeDiff)
		weightedSum += prev.Value * weight
		totalWeight += weight
		prev = next
	}

	if totalWeight == 0 {
		return 0, fmt.Errorf("total weight is zero, cannot compute TWAP")
	}

	return weightedSum / totalWeight, nil
}
