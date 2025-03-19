package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blocky/basm-go-sdk"
)

type DHLTrackingInfo struct {
	Shipments []struct {
		Id     string `json:"id"`
		Status struct {
			Timestamp string `json:"timestamp"`
			Location  struct {
				Address struct {
					CountryCode     string `json:"countryCode"`
					PostalCode      string `json:"postalCode"`
					AddressLocality string `json:"addressLocality"`
				} `json:"address"`
			} `json:"location"`
			StatusCode  string `json:"statusCode"`
			Status      string `json:"status"`
			Description string `json:"description"`
		} `json:"status"`
	} `json:"shipments"`
}

type TrackingInfo struct {
	TrackingNumber string `json:"tracking_number"`
	Address        struct {
		CountryCode     string `json:"countryCode"`
		PostalCode      string `json:"postalCode"`
		AddressLocality string `json:"addressLocality"`
	} `json:"address"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

func getTrackingInfoFromDHL(trackingNumber string, apiKey string) (TrackingInfo, error) {
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL: fmt.Sprintf(
			"https://api-test.dhl.com/track/shipments?trackingNumber=%s",
			trackingNumber,
		),
		Headers: map[string][]string{
			"DHL-API-Key": []string{apiKey},
		},
	}
	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return TrackingInfo{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return TrackingInfo{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	dhlTrackingInfo := DHLTrackingInfo{}
	err = json.Unmarshal(resp.Body, &dhlTrackingInfo)
	if err != nil {
		return TrackingInfo{}, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}

	trackingInfo := TrackingInfo{
		TrackingNumber: dhlTrackingInfo.Shipments[0].Id,
		Address:        dhlTrackingInfo.Shipments[0].Status.Location.Address,
		Status:         dhlTrackingInfo.Shipments[0].Status.Status,
		Description:    dhlTrackingInfo.Shipments[0].Status.Description,
		Timestamp:      dhlTrackingInfo.Shipments[0].Status.Timestamp,
	}

	return trackingInfo, nil
}

type Args struct {
	TrackingNumber string `json:"tracking_number"`
}

type SecretArgs struct {
	DHLAPIKey string `json:"api_key"`
}

//export trackingFunc
func trackingFunc(inputPtr uint64, secretPtr uint64) uint64 {
	var input Args
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return WriteError(outErr)
	}

	var secret SecretArgs
	secretData := basm.ReadFromHost(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	trackingInfo, err := getTrackingInfoFromDHL(
		input.TrackingNumber,
		secret.DHLAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting DHL tracking info: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(trackingInfo)
}

func main() {}
