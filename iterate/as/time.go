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
		DateTime string `json:"dateTime"`
	}{}

	err = json.Unmarshal(resp.Body, &timeData)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"unmarshaling time data: %w...%s",
			err,
			resp.Body,
		)
	}

	parsedTime, err := time.Parse(
		"2006-01-02T15:04:05.0000000",
		timeData.DateTime,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"parsing time data: %w...%s",
			err,
			timeData.DateTime,
		)
	}

	return parsedTime, nil
}
