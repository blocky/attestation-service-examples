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
			"unmarshaling time data: %w...%v",
			err,
			resp.Body,
		)
	}

	// convert datetime from ISO 8601 format to RFC 3339 format then parse
	parsedTime, err := time.Parse(time.RFC3339, timeData.DateTime+"Z")

	return parsedTime, nil
}
