# Random - Generating Random Numbers

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest simple function call that generates random numbers to simulate a die
roll.

Before starting this example, make sure you are familiar with the
[Attesting a Function Call](../attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and the
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

You'll learn how to:

- Generate random numbers in the Blocky AS runtime
- Create a function that returns a random number within a given range

*Note: Random number generation only works on Blocky AS versions 0.1.0-beta.5
and above. While this example may run on previous versions it will generate
deterministic number streams.*

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
  "Value": 3
}
```

## Walkthrough

### Step 1: Create a parameterized random function

We'll demonstrate random number generation by emulating a die roll in `rollDie` in
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the
[`fn-call.json`](./fn-call.json) file contents:

```json
{
  "code_file": "tmp/x.wasm",
  "function": "rollDie",
  "input": {
    "die_sides": 20
  }
}
```

Notice the `input` section, which contains the parameter for `rollDie`,
specifically the `sides` field set to `20`. Meaning, `rollDie` will return
a number in the range [1,20].

Next, we define the `rollDie` function in [`main.go`](./main.go):

```go
type Args struct {
    DieSides int `json:"die_sides"`
}

//export rollDie
func rollDie(inputPtr uint64, secretPtr uint64) uint64 {
    var input Args
    inputData := basm.ReadFromHost(inputPtr)
    err := json.Unmarshal(inputData, &input)
    switch {
    case err != nil:
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
        return WriteError(outErr)
    case input.DieSides < 1:
        outErr := fmt.Errorf("die must have one or more sides")
        return WriteError(outErr)
    }

    roll := rand.Intn(input.DieSides) + 1
    return WriteOutput(roll)
}
```

You will notice a few things:

- First, we fetch the input data and unmarshal it into the
  `Args` struct. We don't fetch any secrets.
- We are able to generate random numbers the way you would in a normal Go
  program (i.e., via `math/rand` or `crypto/rand`). This is because `tinygo`
  provides implementations for `math/rand` and `crypto/rand`, and because we
  properly seed the `wazero` runtime with a source of entropy (`wazero`
  provides a deterministic source of entropy by default). Specifically, we
  seed the runtime with entropy obtained from the
  [AWS Nitro Secure Module (NSM)](https://github.com/aws/aws-nitro-enclaves-nsm-api).

### Step 2: Compile the function to WebAssembly (WASM)

To invoke our function in the Blocky AS server, we first need to compile
it into a WASM file. If you inspect the `build` target in the
[`Makefile`](./Makefile), you'll see the build command:

```bash
@bky-c build . tmp/x.wasm
```

where we use the `Blocky Compiler` (`bky-c`) to compile
[`main.go`](./main.go) to WASM and save it to `tmp/x.wasm`. You can build our
function by calling:

```bash
make build
```

### Step 3: Invoke the function on the Blocky AS server

To invoke the `rollDie` example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response:

```
Output:
{
  "Success": true,
  "Error": "",
  "Value": 3
}
```

You will most likely see a different number for the `Value` field as it is
randomly generated. On that note, subsequent invocations of `rollDie` should
produce different numbers. You can control the range of said numbers by
updating the `sides` field in the `fn-call.json` file.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. Try updating the example to return the result of multiple
rolls.

Check out other examples in this repository, to learn what
else you can do with Blocky AS.
