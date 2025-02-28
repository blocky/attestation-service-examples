# Hello World - Attesting a Function Call

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
  
## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
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

```
Output:
{
  "Success": true,
  "Value": {
    "Number": 42
  }
}
```

## Walkthrough

### Step 1: Create success and error functions

Before we get to error handling, let's create a function that successfully
returns structured data. In [`main.go`](./main.go), we define the
`successFunc` function:

```go
func successFunc(inputPtr, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{42}
	return writeOutput(output)
}
```

This function creates an `Output` struct with a single field, `Number`, and
sets it to 42. It then calls the `writeOutput` function to return it to the user
using the `writeOutput` function.

Let's say, however, that you want to write a function that fail depending on its
starting conditions, or while processing the data it fetches from the web. For
the purpose of this example, we define the `errorFunc` in [`main.go`](./main.go)
that will always fail:

```go
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return writeError(err)
}
```

Both of these functions will successfully complete their execution on the 
Blocky AS server and produce an attestation over their output. What we'd like
is an easy way to tell whether the function succeeded or failed due to a runtime
error.

### Step 2: Use the result pattern

To do this, we can use the [result pattern](https://en.wikipedia.org/wiki/Result_type).
In [`main.go`](./main.go), we define a `Result` struct that can hold either
a successful result or an error:

```go
type Result struct {
	Success bool
	Value   any
}
```

The `Success` field indicates whether the function succeeded or failed, and the
`Value` field holds the result of the function if it succeeded, or an error
string if it failed.

To return a `Result` to user, we need to serialize to bytes and send them to
the `as.WriteToHost` function. Let's say that we want to use JSON to serialize
the `Result` struct. 

We can define the `writeOutput` function to take in our function output, as 
`any`, put it in a `Result` struct, serialize it, and send it to 
`as.WriteToHost`:

```go
func writeOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		as.Log(fmt.Sprintf("Error marshalling result: %s", err))
		return writeError(err)
	}
	return as.WriteToHost(data)
}
```

As you see, we have a challenge here in that the `json.Marshal` function itself
can fail. In this case, we can use the `writeError` function to report that 
error. But wouldn't we run into the same, chicken and egg problem if we
encountered an error in the `writeError` function? Let's take a look.

Our `writeError` function is defined as:

```go
func writeError(err error) uint64 {
	data := Result{}.jsonMarshalWithError(err)
	return as.WriteToHost(data)
}
```

and uses the receiver function `jsonMarshalWithError`:

```go
func (r Result) jsonMarshalWithError(err error) []byte {
	resultStr := fmt.Sprintf(`{ "Success": false, "Value": "%s" }`, err)
	data := []byte(resultStr)
	return data
}
```

to JSON serialize the `Result` struct with the error string. Because we
hand-roll the JSON serialization, we no longer have to worry about the
serialization failing.

### Step 3: Run the functions

To run the `successFunc` function, we can call:

```bash
make run-success
```

which will produce:

```
Output:
{
  "Success": true,
  "Value": {
    "number": 42
  }
}
```

You can see that the `Success` field is set to `true`, which means we can
read the `Value` field as the JSON serialized `Output` struct.

To run the `errorFunc` function, we can call:

```bash
make run-error
```

which will produce:

```
Output:
{
  "Success": false,
  "Value": "expected error"
}
```

Now the `Success` field is set to `false`, which means that we can read the
`Value` field as the error string.

## Next Steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. 

If you explored the 
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example, you can use `Result.Success` to decide whether you want to bring the
attestation on chain or not. You can also extend the example, so that
transactions calling the `User.sol` `verifyAttestedFnCallClaims` function revert
if the `Result.Success` field is set to `false`.

You can also check out other examples in this repository, to learn what
else you can do with Blocky AS.
