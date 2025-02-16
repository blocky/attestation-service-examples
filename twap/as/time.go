package as

import (
	"encoding/json"
	"fmt"
	"time"
)

// TimeNow fetches the current UTC time from an external API.
// In the future this will be implemented as a host function.
func TimeNow() (time.Time, error) {
	req := HostHTTPRequestInput{
		Method: "GET",
		URL:    "https://timeapi.io/api/time/current/zone?timeZone=UTC",
	}
	resp, err := HostFuncHTTPRequest(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("making http request: %w", err)
	}

	timeData := struct {
		Year         int `json:"year"`
		Month        int `json:"month"`
		Day          int `json:"day"`
		Hour         int `json:"hour"`
		Minute       int `json:"minute"`
		Seconds      int `json:"seconds"`
		Milliseconds int `json:"milliSeconds"`
	}{}

	err = json.Unmarshal(resp.Body, &timeData)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"unmarshaling time data: %w...%s",
			err,
			resp.Body,
		)
	}

	parsedTime := time.Date(
		timeData.Year,
		time.Month(timeData.Month),
		timeData.Day,
		timeData.Hour,
		timeData.Minute,
		timeData.Seconds,
		timeData.Milliseconds*1e6,
		time.UTC,
	)

	return parsedTime, nil
}
