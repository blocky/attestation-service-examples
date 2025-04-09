# Accessing Wall Clock Time

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest simple function call that gets the current wall clock time.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and the
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

You'll learn how to:

- Get the current wall clock time

*Note: Wall time only works on Blocky AS versions 0.1.0-beta.6 and above.
While this example may run on previous versions it will generate time based
on a hardcoded start date.*

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
make run
```

You will see the following output extracted from a Blocky AS response:

```
Output:
{
  "Success": true,
  "Error": "",
  "Value": "2025-04-03T20:13:12.9481826Z"
}
```

## Walkthrough

### Step 1: Create a function for getting the current time

By default, the `wazero` runtime does not provide sandboxed functions access to
the system clock. Instead, they are given a hardcoded time that monotonically 
increases by "1ns" for every subsequent call. In the Blocky AS environment,
however, we have configured `wazero` to grant functions access to the system
clock.

This means that your Blocky AS functions can fetch time using standard library
calls same as if they were executing in a normal environment. Indeed, the time
function defined in [`main.go`](./main.go) is a simple wrapper for the
`time.Now()` standard library call.

```go
//export timeNow
func timeNow(_ uint64, _ uint64) uint64 {
    return WriteOutput(time.Now())
}
```

While time related functions are supported, it is important to note that system
sleep functions **are not**. Attempting to use functions that rely on system
sleep (e.g., `time.Sleep()`) will cause your function to panic.

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

To invoke the `timeNow` example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response:

```
Output:
{
  "Success": true,
  "Error": "",
  "Value": "2025-04-03T20:13:12.9481826Z"
}
```

You should see the current wall clock time reflected in the `Value` field.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs.

Check out other examples in this repository, to learn what
else you can do with Blocky AS.
