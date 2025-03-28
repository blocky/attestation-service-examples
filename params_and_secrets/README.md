# Passing parameters and secrets

This example shows you how to pass input parameters and secrets to
functions executed in the Blocky Attestation Service (Blocky AS).

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md)
and the
[Error Handling](../error_handling/README.md)
examples.

You'll learn how to:

- Pass input parameters to functions
- Pass secrets to functions

Every function executed in the Blocky Attestation Service (Blocky AS)
can accept user input. 

There are two function input types:
* Regular input parameters – Sent directly to the Blocky AS server and
passed to your function as-is. These parameters are meant for data that 
does not require any special handling or security measures.
* Secrets – Encrypted before being sent to the Blocky AS Server and only
decrypted inside the server. This protects the integrity of sensitive data
as the key needed for decryption exists exclusively inside the
Blocky AS Server. Secrets are designed for passing sensitive information
such as personal data, API keys, or passwords.

You'll learn how to pass input parameters and secrets by writing a function
that extracts password protected user data. We will pass the data using 
input parameters and the confidential access password using secrets.

* Note: Used [password protection mechanism](./extract.go) is very simple
and cannot be considered secure.
It is meant for demonstration purposes only and is not suitable for any use
outside of this example. *
  
## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/v0.1.0-beta.4/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.

## Quick Start

To run this example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "extracted_data": "your extracted information"
  }
}
```

## Walkthrough

### Step 1: Create a parametrized extraction function

We'll demonstrate both input parameter and secret passing in the `extract`
function in [`main.go`](./main.go). As in previous examples,
we will call this function using the `bky-as` CLI by passing in the
[`fn-call.json`](./fn-call.json) file contents:

```json
{
  "code_file": "tmp/x.wasm",
  "function": "extract",
  "input": {
    "data": "z9s5DqczsNfnKhw86T6EiH74bNaoGS+T+kWNgDf7lvyRqTOyLiZrhloRnVq/nSXuIMD37aE1"
  },
  "secret": {
    "password": "secret-password"
  }
}
```

Input parameters and secrets can be passed into your function using
the `input` and `secret` sections respectively.

Our payload contains the `input` parameter section with
the password protected string inside the `data` field. 
Confidential password used for data extraction is stored 
in the `password` field of the `secret` section.

Note that both `input` and `secret` sections can hold any `json` document, 
and it is up to you, how you wish to structure them.

Next, we define the `extract` function in [`main.go`](./main.go):

```go
type Parms struct {
    Data string `json:"data"`
}

type Secrets struct {
    Password string `json:"password"`
}

type Output struct {
    ExtractedData string `json:"extracted_data"`
}

//export extract
func extract(inputPtr uint64, secretPtr uint64) uint64 {
	var params Parms
    inputData := basm.ReadFromHost(inputPtr)
    err := json.Unmarshal(inputData, &params)
    if err != nil {
        outErr := fmt.Errorf("unmarshaling params: %w", err)
	    return WriteError(outErr)
    }
  
    var secrets Secrets
    secretData := basm.ReadFromHost(secretPtr)
    err = json.Unmarshal(secretData, &secrets)
    if err != nil {
        outErr := fmt.Errorf("unmarshaling secrets: %w", err)
        return WriteError(outErr)
    }
  
    result, err := extractData(params.Data, secrets.Password)
    if err != nil {
        outErr := fmt.Errorf("extracting data: %w", err)
        return WriteError(outErr)
    }
  
    output := Output{ExtractedData: result}
  
    return WriteOutput(output)
}
```

You will notice a few things:

- The function takes two `uint64` arguments and returns a `uint64`. These are
  fat pointers to shared memory managed by the Blocky AS server, where the first
  32 bits are a memory address and the second 32 bits are the size of the data.
  The memory space is sandboxed and shared between the TEE host program (Blocky
  AS server) and the WASM runtime (your function). 
- The `inputPtr` and `secretPtr` arguments carry user-defined function
  input and secrets. They are used to access the data provided in the `input`
  and `secret` sections of the `bky-as` [CLI input](./fn-call.json)
- The function uses the `basm`
  [Blocky Attestation Service WASM Go SDK](https://github.com/blocky/basm-go-sdk)
  `basm.ReadFromHost` to fetch the input parameter and secret data
  and unmarshal then into the `Params` and `Secrets` structs. These structs need
  to match the `json` structure of the `input` and `secret` sections of the 
  [CLI input](./fn-call.json).
- The predefined [`extractData`](./extract.go) function accepts the `params.Data`
  password protected string and the `secrets.Password` that allows for its safe
  extraction. It returns the extracted string value and passes it to the function
  output.
- The details of the [`extractData`](./extract.go) function is out of this
  example's scope.
- The output of the function is also a memory pointer, whose value will be
  returned to the user.



### Step 2: Compile the function to WebAssembly (WASM)

To invoke our function in the Blocky AS server, we first need to compile
it into a WASM file. If you inspect the `build` target in the
[`Makefile`](./Makefile), you'll see the build command:

```bash
@docker run --rm \
    -v .:/src \
    -w /src \
    tinygo/tinygo:0.31.2 \
    tinygo build -o tmp/x.wasm -target=wasi ./...
```

where we use `docker` to run [`tinygo`](https://tinygo.org/) to compile
[`main.go`](./main.go) to WASM and save it to `tmp/x.wasm`. You can build our
function by calling:

```bash
make build
```

### Step 3: Invoke the function on the Blocky AS server

To invoke the `extract` example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "extracted_data": "your extracted information"
  }
}
```

Note that in `bky-as` CLI [output](./tmp/successout.json) you can also find 
the `hash_of_input` and `hash_of_secrets` fields that contain the `SHA3/512`
hashes of the user provided input and secrets.

If you wish to attempt running the `extract` function with incorrect password 
and then observe the error, call:

```bash
make run-error
```

You will see the following output extracted from a Blocky AS response:

```json
{
  "Success": false,
  "Error": "could not extract data: incorrect password",
  "Value": null
}
```

## Next Steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. Try modifying the example to extract multiple password
protected strings in one call. 
You can use the predefined input from
[`fn-call.json`](./fn-call.json) multiple times.

Check out other examples in this repository, to learn what
else you can do with Blocky AS.
