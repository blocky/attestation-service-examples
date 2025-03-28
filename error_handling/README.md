# Error Handling

This example shows you a useful pattern for reporting errors during the 
execution of functions in the Blocky Attestation Service (Blocky AS).

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md)
example.

In this example, you'll learn how to:

- Use the [result pattern](https://en.wikipedia.org/wiki/Result_type) in your
functions
- Return structured data from your functions
- Report errors from your functions
- Log errors in your functions
  
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
make run-success
```

You will see the following output extracted from a Blocky AS response:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "number": 42
  }
}
```

## Walkthrough

### Step 1: Create success and error functions

Before we get to error handling, let's create a function that successfully
returns structured data. In [`main.go`](./main.go), we define the
`successFunc` function:

```go
func successFunc(inputPtr uint64, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{Number: 42}
	return WriteOutput(output)
}
```

This function creates an `Output` struct with a single field, `Number`, and
sets it to `42`. It then calls the `WriteOutput` function to return it to the
user using the `WriteOutput` function.

Let's say, however, that you want to write a function that may fail depending on
its starting conditions, or while processing the data it fetches from the web.
For the purpose of this example, we define the `errorFunc` in 
[`main.go`](./main.go) that will always fail:

```go
func errorFunc(inputPtr uint64, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return WriteError(err)
}
```

Both `successFunc` and `errorFunc` will successfully complete their execution on
the Blocky AS server and produce an attestation over their output. What we'd
like is an easy way to tell whether the function succeeded or failed due to a
runtime error.

### Step 2: Use the result pattern

To do this, we can use the [result pattern](https://en.wikipedia.org/wiki/Result_type).
In [`output.go`](./output.go), we define a `Result` struct that can hold either
a successful result or an error:

```go
type Result struct {
	Success bool
	Error   string
	Value   any
}
```

The `Success` field indicates whether the function succeeded or failed. The
`Error` field holds an error message if the function failed (`Success` is
`false`). If the function succeeds (`Success` is `true`), the `Error` field
should be disregarded. The `Value` field holds the result of the function if it
succeeded (`Success` is `true`). If the function fails (`Success` is `false`),
the `Value` field should be disregarded. Note, that this three-field struct
design of `Result` is a slight departure from the traditional result pattern, 
where typically a `Value` field of type `any` might hold either an error message, or another return
type. Having the three fields, however, allows for one pass parsing of `Result`
structs in the client code.

To return a `Result` to user, we need to serialize to bytes and send them to the
`basm`
[Blocky Attestation Service WASM Go SDK](https://github.com/blocky/basm-go-sdk/tree/v0.1.0-beta.5)
`basm.WriteToHost` function. Let's say that we want to use JSON to serialize the
`Result` struct.

We can define the `WriteOutput` function in [`output.go`](./output.go) to take
in our function output, as `any`, put it in a `Result` struct, serialize it, and
send it to `basm.WriteToHost`:

```go
func WriteOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		basm.Log(fmt.Sprintf("Error marshalling Result: %v", err))
		return WriteError(err)
	}
	return basm.WriteToHost(data)
}
```

As you see, we have a challenge here in that the `json.Marshal` function itself
can fail. In this case, we can use the `WriteError` function to report that 
error. But wouldn't we run into the same, chicken-and-egg problem if we
encountered an error in the `WriteError` function? Let's take a look.

Our `WriteError` function in [`output.go`](./output.go) is defined as:

```go
func WriteError(err error) uint64 {
	data := Result{}.JSONMarshalWithError(err)
	return basm.WriteToHost(data)
}
```

and uses the receiver function `JSONMarshalWithError` in:

```go
func (r Result) JSONMarshalWithError(err error) []byte {
	if err == nil {
		err = errors.New("JSONMarshalWithError invoked with nil error")
	}
	resultStr := fmt.Sprintf(
		`{ "Success": false, "Error": "%s" , "Value": null }`,
		err.Error(),
	)
	return []byte(resultStr)
}
```

to JSON serialize the `Result` struct with an error. Because we
hand-roll the JSON serialization, we no longer have to worry about the
serialization failing.

### Step 3: Run the functions

To run the `successFunc` in [`main.go](./main.go) function, we can call:

```bash
make run-success
```

which will produce:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "number": 42
  }
}
```

You can see that the `Success` field is set to `true`, which means we can
disregard the `Error` field and read the `Value` field as the JSON serialized
`Output` struct.

To run the `errorFunc` function, we can call:

```bash
make run-error
```

which will produce:

```json
{
  "Success": false,
  "Error": "expected error",
  "Value": null
}
```

Now the `Success` field is set to `false`, which means that we can read the
error message from the`Error` field and disregard the `Value` field.

### Step 4: Handle runtime errors

But what if the function encounters a runtime error? In
[`main.go`](./main.go), we define the `panicFunc` function:

```go
func panicFunc(inputPtr uint64, secretPtr uint64) uint64 {
	basm.LogToHost("Expected panic call\n")
	panic(nil)
}
```

where we call `panic` to trigger a runtime error during the function execution.
To debug code with runtime errors, we might want to emit information on the
state of the program before the panic. To do that, we can use the
`basm.LogToHost` function from the
[Blocky Attestation Service WASM Go SDK](https://github.com/blocky/basm-go-sdk).
When you run the function on a local instance of the Blocky AS server (the
`host` value in `config.toml` is set to `"local-server"`) you will be able to
see the log messages sent to host.

To run the `panicFunc` function, we can call:

```bash
make run-panic
```

which will produce output similar to:

```
🚀 Starting local server at http://127.0.0.1:8081 ...success
Expected panic call
2025/03/19 17:41:20 Error: making fn call attestation request: making function call: unexpected status, wanted '200' got '422': "invoking code: error invoking code 'panicFunc': calling function: wasm error: unreachable\nwasm stack trace:\n\t.runtime._panic(i32,i32)\n\t\t0x3fa1: /usr/local/tinygo/src/runtime/runtime_tinygowasm.go:70:6\n\t.panicFunc(i64,i64) i64\n\t\t0x4a7a0: /src/main.go:27:7 (inlined)"
```

The error message is quite verbose, but it tells us that the `panicFunc`
function encountered a runtime error in `main.go` at line 27. We can also see
the context for the error from out `Expected panic call` log message shown above
the server error message.

## Next Steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. 

If you explored the 
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example, you can extend it to use `Result.Success` to decide whether you want to
bring the attestation on chain or not. You can also extend the example, so that
transactions calling the `User.sol` `verifyAttestedFnCallClaims` function revert
if the `Result.Success` field is set to `false`.

You can also check out other examples in this repository, to learn what
else you can do with Blocky AS.
