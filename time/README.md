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

- Access the current wall clock time in your WASM function

*Note: Wall time only works on Blocky AS versions 0.1.0-beta.6 and above.
While this example may run on previous versions it will generate time based
on a hardcoded start date.*

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/{{AS_VERSION}}/setup)
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

Normally, the Blocky AS runtime does not provide sandboxed functions access to
the system clock. Instead, they are given a hardcoded time that monotonically 
increases by "1ns" for every subsequent call. 

Starting with 0.1.0-beta.6, the Blocky AS runtime fetches the current time
from a remote time-server and provides these values to functions in the
sandboxed environment via standard system calls.

Starting with 0.1.0-beta.9, the Blocky AS runtime fetches the current time
directly from AWS Nitro Enclave's Hypervisor via the
[`PTP`](https://en.wikipedia.org/wiki/Precision_Time_Protocol) device
(`/dev/ptp0`) and provides these values to functions in the sandboxed
environment via standard system calls. 

Getting the time from `/dev/ptp0` device means asking a dedicated clock
provided by the AWS Nitro Hypervisor for the current timestamp.
Instead of using a standard system clock, it gives you direct access to
a precise and reliable time source maintained by the AWS cloud infrastructure.
Each time you request the current time, the Blocky AS runtime
talks to the Nitro hypervisor through the `/dev/ptp0` device and returns
an accurate timestamp.
You can read more about time in AWS Nitro Enclaves in this
[article](https://evervault.com/blog/how-we-built-enclaves-resolving-clock-drift-in-nitro-enclaves).

This means that your Blocky AS functions can fetch time using standard library
calls, like `time.Now()`, same as if they were executing
in a non-TEE environment.

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
