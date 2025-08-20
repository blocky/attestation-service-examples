# Accessing Wall Clock Time

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest simple function call that gets the current wall clock time.

Before starting this example, make sure you are familiar with the
[Attesting a Function Call](../attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and the
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

You'll learn how to:

- Access the current wall clock time in your WASM function

*Note: Wall time only works on Blocky AS versions 0.1.0-beta.6 and above.
While this example may run on previous versions it will generate time based
on a hardcoded start date.*

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://docs.blocky.rocks/attestation-service/{{{AS_VERSION}}}/setup)
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

Before version 0.1.0.beta-9, the Blocky AS runtime does not provide sandboxed
functions access to the system clock. Instead, they are given a hardcoded
time that monotonically increases by "1ns" for every subsequent call.

Starting with 0.1.0-beta.9, guest functions have access
to time from [`PTP`](https://en.wikipedia.org/wiki/Precision_Time_Protocol).
Each time a guest function calls a standard library time function,
like `time.Now()`,  the resulting time system call triggers
the Blocky AS runtime to fetch current time from the `/dev/ptp0` device.
Read more about this design in the
[article](https://evervault.com/blog/how-we-built-enclaves-resolving-clock-drift-in-nitro-enclaves).

We define a simple function in [`main.go`](./main.go) that calls `time.Now()`.

```go
//export timeNow
func timeNow(_ uint64, _ uint64) uint64 {
    return WriteOutput(time.Now())
}
```

While time related functions are supported, it is important to note that system
sleep functions **are not**. Attempting to use functions that rely on system
sleep (e.g., `time.Sleep()`) will cause your function to panic.

### Step 2: Invoke the function on the Blocky AS server

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
