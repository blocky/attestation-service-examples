# Tracking a Shipment with DHL

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches and processes data from the DHL Tracking API.

Before starting this example, make sure you are familiar with the
[Attesting a Function Call](../attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md)
and the
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

In this example, you'll learn how to:

- Pass in parameters and secrets to your function
- Make an HTTP request to an external API in your function
- Parse a JSON response from an API

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://docs.blocky.rocks/attestation-service/v0.1.0-beta.12/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.

## Quick Start

To run this example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response showing
you the status of a DHL shipment with tracking number `00340434292135100186`.

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "tracking_number": "00340434292135100186",
    "address": {
      "countryCode": "US",
      "postalCode": "89014",
      "addressLocality": "Henderson, NV, US"
    },
    "status": "DELIVERED",
    "description": "DELIVERED - PARCEL LOCKER",
    "timestamp": "2023-05-08T10:37:00"
  }
}
```

> Note that this demo uses a demo DHL API key. If you run up against rate 
> limits, you can get your own API key by signing up for a 
> [DHL developer account](https://developer.dhl.com/) and updating the
> `api_key` in [`fn-call.json`](./fn-call.json) with your DHL developer API key.

## Walkthrough

Let's say you want to implement an oracle that fetches the status of a DHL
shipment using the DHL API:

```bash
curl -s \
    'https://api-test.dhl.com/track/shipments?trackingNumber=00340434292135100186' \
    -H 'DHL-API-Key: demo-key' \
  | jq .
```

If you run the above command, you will get a lot of information. Let's say that 
you want to parse out just a few details about the shipment relevant to your 
application.

### Step 1: Create a parameterized oracle function

We'll implement the oracle as `trackingFunc` in
[`main.go`](./main.go). We will call this function using the `bky-as` CLI by
passing in the [`fn-call.json`](./fn-call.json) file contents:

```json
{
  "code_file": "tmp/x.wasm",
  "function": "trackingFunc",
  "input": {
    "tracking_number": "00340434292135100186"
  },
  "secret": {
    "api_key": "demo-key"
  }
}
```

Notice the `input` section, which contains the `tracking_number` parameter, and
the `secret` section, which  contains a DHL `api_key`.

Next, we define the `trackingFunc` function in [`main.go`](./main.go):

```go
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
```

First, we get the input parameters and secrets. Next, we call
the `getTrackingInfoFromDHL` function to fetch tracking information for 
`input.TrackingNumber` using the `secret.DHLAPIKey` API key. Finally, we 
return the `trackingInfo` to user by converting its data to a fat pointer
using the `WriteOutput` function and returning the pointer from `trackingFunc`
to the Blocky AS server host runtime.

### Step 2: Make a request to the DHL API

The `getTrackingInfoFromDHL` function, called by `trackingFunc`, will make an
HTTP request to the DHL API to fetch the tracking information for a specific
tracking number.

Let's start by setting up a struct in [`main.go`](./main.go) to parse the
relevant fields from the DHL API response JSON:

```go
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
```

Next, we'll define the `getTrackingInfoFromDHL` function to fetch and parse the
data from the DHL API:

```go
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

func getTrackingInfoFromDHL(
	trackingNumber string,
	apiKey string,
) (
	TrackingInfo,
	error,
) {
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
```

The `getTrackingInfoFromDHL` function takes the `trackingNumber`, and `apiKey`
as arguments. First it constructs an HTTP request to the DHL API using
`trackingNumber` in the URL and the `apiKey` in the headers. It then sends the
request to the `basm.HTTPRequest` function, which makes the request through the
Blocky AS server networking stack. Next, it checks the response status code and
unmarshalls the JSON response into the `DHLTrackingInfo` struct. Finally, it
parses out information from a `DHLTrackingInfo` struct into a `TrackingInfo`
struct and returns it to the `trackingFunc`. The `trackingFunc` function returns
a `Result` containing the `TrackingInfo` to the Blocky AS server to create an
attestation over the function call and the `Result` struct.

### Step 3: Run the oracle

To run `trackingFunc`, you need call:

```bash
make run
```

You'll see output similar to the following:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "tracking_number": "00340434292135100186",
    "address": {
      "countryCode": "US",
      "postalCode": "89014",
      "addressLocality": "Henderson, NV, US"
    },
    "status": "DELIVERED",
    "description": "DELIVERED - PARCEL LOCKER",
    "timestamp": "2023-05-08T10:37:00"
  }
}
```

where `"Success": true,` tells you that the function call was successful and 
the `Value` field gives you a JSON-serialized `TrackingInfo` struct.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For example, you can try passing in different tracking
numbers to `trackingFunc`, or even multiple tracking numbers with some
modifications to the code. You can also change the API endpoint in
`getTrackingInfoFromDHL` to fetch data from a different API, or even multiple
APIs. You may also want to explore the
[Bringing A Blocky AS Function Call Attestation On Chain](../on_chain/README.md)
example to learn you can bring the `TrackingInfo` struct into a smart contract.
