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

	// Remove duplicate samples with the same timestamp
	uniqueSamples := make([]Price, 0, len(samples))
	seen := make(map[time.Time]struct{})
	for _, sample := range samples {
		if _, ok := seen[sample.Timestamp]; !ok {
			uniqueSamples = append(uniqueSamples, sample)
			seen[sample.Timestamp] = struct{}{}
		}
	}

	if len(uniqueSamples) == 1 {
		return uniqueSamples[0].Value, nil
	}

	// Sort samples from latest to earliest
	lessThan := func(i, j int) bool {
		return uniqueSamples[i].Timestamp.After(uniqueSamples[j].Timestamp)
	}
	sort.Slice(uniqueSamples, lessThan)

	var weightedSum, totalWeight float64

	// IMPORTANT: The value of the last sample is not included in the calculation
	// because it doesn't have a next sample to compare with. However, its
	// timestamp is used to calculate the weight of the previous sample.
	prev := uniqueSamples[0]
	for _, next := range uniqueSamples[1:] {
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
