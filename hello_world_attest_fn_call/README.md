# Hello World

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a simple function calls.

You'll learn how to:

- Create a function that returns a `"Hello, World!"` message
- Create a function that writes a message to a log and returns a
  user-defined error
- Invoke functions in the Blocky AS using its `bky-as` CLI
- Extract function output from the Blocky AS attestation

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.

## Quick Start

To run this example, call the following command:

```bash
make hello-world
```

You will see the following JSON output extracted from a Blocky AS attestation:

```
Result:
{
  "Success": true,
  "Error": "",
  "Output": "Hello, World!"
}
```

where `Success` indicates whether the function call was successful, `Error`
contains any error messages, and `Output` contains the output of the function,
in this case the string "Hello, World!".

## Walkthrough

__Step 1: Create a function that returns a "Hello, World!" message.__

Our first goal is to create a simple function that returns a `"Hello, World!"`
message. We will write this function in Go and compile it to WebAssembly (WASM)
to run on the Blocky AS server. If you open [`main.go`](./main.go) you'll see
our function there:

```go
//export helloWorld
func helloWorld(inputPtr, secretPtr uint64) uint64 {
    return as.WriteOutput("Hello, World!")
}
```

You will notice a few things:

- The `helloWorld` function is exported so that it can be invoked by the
  Blocky AS server in a WASM runtime.
- The function takes two `uint64` arguments and returns a `uint64`. These are
  fat pointers to shared memory managed by the Blocky AS server, where the first
  32 bits are a memory address and the second 32 bits are the size of the data.
  The memory space is sandboxed and shared between the TEE host program (Blocky
  AS server) and the WASM runtime (your function). The `inputPtr` and
  `secretPtr` arguments used to pass user-defined function input and secrets,
  through we don't make use of them in this example. The output of the function
  is also a pointer to memory, whose value will be returned to the user.
- The function calls the `as.WriteOutput` function to write the string
  `"Hello, World!"` to the Blocky AS server's output buffer. This function
  returns a pointer to the output buffer, which is then returned to the user. We
  provide the functions in the `as` package as a part of our SDK.

__Step 2: Create a function that writes a message to a log and returns a
user-defined error__

In [`main.go`](./main.go), you'll also see a function called `helloError`:

```go
//export helloError
func helloError(inputPtr, secretPtr uint64) uint64 {
    as.Log("Returning an expected error")
    return as.WriteError("expected error")
}
```

You will notice a call to `as.Log` to write a message to the Blocky AS server's
log. This is useful for debugging and monitoring the function's behavior.
You'll also see a call to `as.WriteError` to return an error message to the
user.

__Step 3: Invoke functions in the Blocky AS using its `bky-as` CLI__

To invoke these functions in the Blocky AS server, we first need to compile
them into an executable WASM file. If you inspect the `build` target in the
[`Makefile`](./Makefile), you'll see the build command:

```bash
@docker run --rm \
  -v .:/src \
  -w /src \
  tinygo/tinygo:0.31.2 \
  tinygo build -o tmp/x.wasm -target=wasi main.go
```

where we use `docker` to run [TinyGo](https://tinygo.org/) to compile 
[`main.go`](./main.go) to WASM. The resulting WASM file is saved to 
`tmp/x.wasm`.

Next, we use the `bky-as` CLI to invoke the functions in the Blocky AS server.
If you inspect the `run` target in the [`Makefile`](./Makefile), you'll see the
command:

```bash
@cat $(FUNCTION)-call.json \
	| bky-as attest-fn-call > tmp/out.json
```

where we use `cat` to read the function call JSON from a file, pipe it to
`bky-as attest-fn-call`, and save the output to `tmp/out.json`.
The input JSON file in this example will be either `hello-world-call.json` or
`hello-error-call.json`, depending on the value of the `FUNCTION` variable.
If you inspect [`hello-world-call.json`](./fn-call.json), you'll see:

```json
[
  {
    "code_file": "tmp/x.wasm",
    "function": "helloWorld"
  }
]
```

where `code_file` is the path to the WASM file we compiled earlier, and
`function` is the name of the exported function we want to call.

To run these functions, you can call:

```bash
make hello-world
```

or

```bash
make hello-error
```

__Step 4: Extract function output from the Blocky AS attestation__

The `run` target will extract the attested output of the function calls.
And so, for:

```bash
make hello-world
```

the output will show the `Result:` section:

```
Result:
{
  "Success": true,
  "Error": "",
  "Output": "Hello, World!"
}
```

where `"Success": true` indicates that the function call was successful,
`"Error": ""` is empty and can be ignored since the call was successful, and
`"Output": "Hello, World!"` contains the expected output of the function call. 

Likewise, when we call:

```bash
make hello-error
```

The output shows:

```
Result:
{
  "Success": false,
  "Error": "expected error",
  "Output": null
}
Logs:
Returning an expected error
```

where the `Result:` sections shows `"Success": false` indicates that the function call was, as expected, unsuccessful, `"Error": "expected error"` contains the expected error message, and `"Output": null`  can be ignored since the call was unsuccessful.
You can also see the `Logs:` section, which shows the log message written by the function.

To dive deeper, let's again look at the `run` target in the [`Makefile`](./Makefile). There you will see that we save the output of the `bky-as attest-fn-call` command to `tmp/out.json`. For the `hello-error` function, 
`tmp/out.json` contains:

```json
{
  "enclave_attested_application_public_key": {
    "enclave_attestation": {
      "Platform": "nitro",
      "PlAttests": [
        "hEShATgioFkRd6lpbW9kdWxlX2lkeCdpLTAyYTIzMTAzYWQ1OGY0YTc4LWVuYzAxOTRmYzJhNDM3YmM5MzNmZGlnZXN0ZlNIQTM4NGl0aW1lc3RhbXAbAAABlTlGHg5kcGNyc7AAWDDt1b/NhX4g/18VH7uBYPVPY47Fvd7lki0sXbNXXJxNsUuIpT2QkXllxwW726SO5lcBWDBLTVs2YbPvwSkgkAyA4Sbkzng8Ui3mwCoqW/evOiuTJ7hndvGI5L4cHEBKEp29pJMCWDDb7K6bLp55sQvPVROKfQKOVl2FzqewlWOQADBOHULW0taDnr7qA9VeFMA+xHSIWuMDWDBKHW5YCApLVDwgqPRTiPUQEQm0kgJHWizZ0SKqe6/QFm773wKTBtlDec+zkitLw+wEWDAMx2xMbZW4RLf3ukTMjsLJjMaOQzX54V5ItB2O46kXTydQcpwffXQtp1+TrWg/nwkFWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAPWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABrY2VydGlmaWNhdGVZAn4wggJ6MIICAaADAgECAhABlPwqQ3vJMwAAAABnvL3OMAoGCCqGSM49BAMDMIGOMQswCQYDVQQGEwJVUzETMBEGA1UECAwKV2FzaGluZ3RvbjEQMA4GA1UEBwwHU2VhdHRsZTEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxOTA3BgNVBAMMMGktMDJhMjMxMDNhZDU4ZjRhNzgudXMtd2VzdC0yLmF3cy5uaXRyby1lbmNsYXZlczAeFw0yNTAyMjQxODQzMjNaFw0yNTAyMjQyMTQzMjZaMIGTMQswCQYDVQQGEwJVUzETMBEGA1UECAwKV2FzaGluZ3RvbjEQMA4GA1UEBwwHU2VhdHRsZTEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxPjA8BgNVBAMMNWktMDJhMjMxMDNhZDU4ZjRhNzgtZW5jMDE5NGZjMmE0MzdiYzkzMy51cy13ZXN0LTIuYXdzMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEkHUTUMrn0Cx/SjmnGQQx+uElzJN5GgoZtk8a4TO/7GM4n2VD/pytpd3fA+yVNVeqySdUGQ7VLr4zqmrAaDbvigk1kS1E9udtb8FIPqGTKDH9MrRoJDQkkPJv+Jcn8Glcox0wGzAMBgNVHRMBAf8EAjAAMAsGA1UdDwQEAwIGwDAKBggqhkjOPQQDAwNnADBkAjB0YPmYUhmaIztSxIJPFrL6PhaMf1GwbdcyObxhl0eJyOEG39S4t8QYdcROaWmD9lMCMF7fwje+yBIOxZhFg2Vf3o40GSdDFzOsSWVs6Id+l9Hm+RW7/XDs8mvwV4ura+SsM2hjYWJ1bmRsZYRZAhUwggIRMIIBlqADAgECAhEA+TF1aBuQr+EdRsy05Of4VjAKBggqhkjOPQQDAzBJMQswCQYDVQQGEwJVUzEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxGzAZBgNVBAMMEmF3cy5uaXRyby1lbmNsYXZlczAeFw0xOTEwMjgxMzI4MDVaFw00OTEwMjgxNDI4MDVaMEkxCzAJBgNVBAYTAlVTMQ8wDQYDVQQKDAZBbWF6b24xDDAKBgNVBAsMA0FXUzEbMBkGA1UEAwwSYXdzLm5pdHJvLWVuY2xhdmVzMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE/AJU66YIwfNocOKa2pC+RjgyknNuiUv/9nLZiURLUFHlNKSx9tvjwLxYGjK3sXYHDt4S1po/6iEbZudSz33R3QlfbxNw9BcIQ9ncEAEh5M9jASgJZkSHyXlihDBNxT/0o0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSQJbUN2QVH55bDlvpync+Zqd9LljAOBgNVHQ8BAf8EBAMCAYYwCgYIKoZIzj0EAwMDaQAwZgIxAKN/L5Ghyb1e57hifBaY0lUDjh8DQ/lbY6lijD05gJVFoR68vy47Vdiu7nG0w9at8wIxAKLzmxYFsnAopd1LoGm1AW5ltPvej+AGHWpTGX+c2vXZQ7xh/CvrA8tv7o0jAvPf9lkCwjCCAr4wggJFoAMCAQICEQCLtgI6xJuqK6j3iacjKfIPMAoGCCqGSM49BAMDMEkxCzAJBgNVBAYTAlVTMQ8wDQYDVQQKDAZBbWF6b24xDDAKBgNVBAsMA0FXUzEbMBkGA1UEAwwSYXdzLm5pdHJvLWVuY2xhdmVzMB4XDTI1MDIyMDIwMjMyMFoXDTI1MDMxMjIxMjMyMFowZDELMAkGA1UEBhMCVVMxDzANBgNVBAoMBkFtYXpvbjEMMAoGA1UECwwDQVdTMTYwNAYDVQQDDC00OGE2ZDFhZTcyZTlhMDczLnVzLXdlc3QtMi5hd3Mubml0cm8tZW5jbGF2ZXMwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAAR6+HeRSGgrxnOwkoxv8OKj2b3GbWqClLl6lLx4X7aJ6iPBJgX3WIkKt7o6qdH3tFaxhPO/TypAgEV/EjC5wpF5uDjSOUkOmCl3Nrj5zYbqFfspB/jmZuI5i/0f59pwrVajgdUwgdIwEgYDVR0TAQH/BAgwBgEB/wIBAjAfBgNVHSMEGDAWgBSQJbUN2QVH55bDlvpync+Zqd9LljAdBgNVHQ4EFgQUPrxPlv/nkqBUfRRgY2wpcACKQuQwDgYDVR0PAQH/BAQDAgGGMGwGA1UdHwRlMGMwYaBfoF2GW2h0dHA6Ly9hd3Mtbml0cm8tZW5jbGF2ZXMtY3JsLnMzLmFtYXpvbmF3cy5jb20vY3JsL2FiNDk2MGNjLTdkNjMtNDJiZC05ZTlmLTU5MzM4Y2I2N2Y4NC5jcmwwCgYIKoZIzj0EAwMDZwAwZAIwDjC93DcOXYeJWV8GHAnsN5x/DXG87IwbaD6flHCdQlAqlrrJDi16bI8xUcVDTGp3AjAq+93d4I45BxfuD4mon8CdcTujlqkcRHqd9qSmANBSv/UX6YUKbGmm9g7Myzooi1FZAxgwggMUMIICm6ADAgECAhEAw0RDGqXRvPYiiI+KlwsUzzAKBggqhkjOPQQDAzBkMQswCQYDVQQGEwJVUzEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxNjA0BgNVBAMMLTQ4YTZkMWFlNzJlOWEwNzMudXMtd2VzdC0yLmF3cy5uaXRyby1lbmNsYXZlczAeFw0yNTAyMjQwMDU1MDJaFw0yNTAzMDIwMTU1MDJaMIGJMTwwOgYDVQQDDDMzNmJiNzhkYzU4YjBhNWY1LnpvbmFsLnVzLXdlc3QtMi5hd3Mubml0cm8tZW5jbGF2ZXMxDDAKBgNVBAsMA0FXUzEPMA0GA1UECgwGQW1hem9uMQswCQYDVQQGEwJVUzELMAkGA1UECAwCV0ExEDAOBgNVBAcMB1NlYXR0bGUwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAASp9EzU6mADkRGK3cIMmaA9o6or6v0yWctk2nAVHgpEqgsCKTo8U/Q7f7UYmWQGOZMGTQ/Yqk246MxLdXUS0fLr+gQASRpy+LK48r7rwEPrWt8yqwliWIp1UNtFeXUDcEujgeowgecwEgYDVR0TAQH/BAgwBgEB/wIBATAfBgNVHSMEGDAWgBQ+vE+W/+eSoFR9FGBjbClwAIpC5DAdBgNVHQ4EFgQUN5HX6Ri+Ac+1PCknMiZfuEXdoVUwDgYDVR0PAQH/BAQDAgGGMIGABgNVHR8EeTB3MHWgc6Bxhm9odHRwOi8vY3JsLXVzLXdlc3QtMi1hd3Mtbml0cm8tZW5jbGF2ZXMuczMudXMtd2VzdC0yLmFtYXpvbmF3cy5jb20vY3JsLzM4MjgzMjFiLTBkNDgtNGZjMy04MTUyLWEzZjQ1YTNiMzBlZi5jcmwwCgYIKoZIzj0EAwMDZwAwZAIwGrRzhyqMrcFOge6buCdSIzfdrIZ+STLumfoC4o5FCpJ4tnMogDTVFHD5lKJ+dgSeAjBJrzaBs/sKWEbzdSjORdlRXui4iGeg2Wqnv0dDebBe5UFbdAXU+Zwo1DQar9wrO2pZAsIwggK+MIICRKADAgECAhQjyijWQmLyfnH8IKM4PjGeCXQ4GjAKBggqhkjOPQQDAzCBiTE8MDoGA1UEAwwzMzZiYjc4ZGM1OGIwYTVmNS56b25hbC51cy13ZXN0LTIuYXdzLm5pdHJvLWVuY2xhdmVzMQwwCgYDVQQLDANBV1MxDzANBgNVBAoMBkFtYXpvbjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAldBMRAwDgYDVQQHDAdTZWF0dGxlMB4XDTI1MDIyNDA5MjMzM1oXDTI1MDIyNTA5MjMzM1owgY4xCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApXYXNoaW5ndG9uMRAwDgYDVQQHDAdTZWF0dGxlMQ8wDQYDVQQKDAZBbWF6b24xDDAKBgNVBAsMA0FXUzE5MDcGA1UEAwwwaS0wMmEyMzEwM2FkNThmNGE3OC51cy13ZXN0LTIuYXdzLm5pdHJvLWVuY2xhdmVzMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE3I1rZjXOUyVx4jHQIrLbj7NNeVLtZTvUPumM0G1zLfD9TFEDofUY4sfRTG+tqeiJe8G8pZ3djhMWmD+kNvir3vlV9efystofcIIFQB7+VuSBn8XHPHOpnqpIDr4KPUh8o2YwZDASBgNVHRMBAf8ECDAGAQH/AgEAMA4GA1UdDwEB/wQEAwICBDAdBgNVHQ4EFgQUVypfioMXTXFrZfqZD/cvaC4+xuUwHwYDVR0jBBgwFoAUN5HX6Ri+Ac+1PCknMiZfuEXdoVUwCgYIKoZIzj0EAwMDaAAwZQIwQU2Q4weVl5wU/sT/8w/eM53UI9jUp0STqHKFPMGrSV3ibv0emBZX/oGdmtA0ly5QAjEA+yb1Q0wnZoyUeP+b5i6LibCCNHsP1+Espx3cEB/GBdMixCy5GrjzY3Rtjm4CzrHdanB1YmxpY19rZXn2aXVzZXJfZGF0YVh5eyJjdXJ2ZV90eXBlIjoicDI1NmsxIiwiZGF0YSI6IkJGTlN1WHBqRHJ2blpPbkl5NmhQakpHczQwb01DWEVpS3NMdDg3Q0lFbXJYL3JqbXBBN0U4K1k1THE3dFFlb0I4SFJ3WTFCTm5TdGNXWUJpc2RYUUhFRT0ifWVub25jZfZYYJa/h4focWY65nh2MFN9N6hdENB2Vzw/04ndd/CkOYmFRNukipbdsZ074Y3K0/MxmJwpbe5WK7RtFt22xkUzXwRiB8uhhewOyTUzPqfwE0aR0sngECoLf8oJ4kGeSFgc7w==",
        "hEShATgioFkRHqlpbW9kdWxlX2lkeCdpLTAyYTIzMTAzYWQ1OGY0YTc4LWVuYzAxOTRmYzJhNDM3YmM5MzNmZGlnZXN0ZlNIQTM4NGl0aW1lc3RhbXAbAAABlTlGHg9kcGNyc7AAWDDt1b/NhX4g/18VH7uBYPVPY47Fvd7lki0sXbNXXJxNsUuIpT2QkXllxwW726SO5lcBWDBLTVs2YbPvwSkgkAyA4Sbkzng8Ui3mwCoqW/evOiuTJ7hndvGI5L4cHEBKEp29pJMCWDDb7K6bLp55sQvPVROKfQKOVl2FzqewlWOQADBOHULW0taDnr7qA9VeFMA+xHSIWuMDWDBKHW5YCApLVDwgqPRTiPUQEQm0kgJHWizZ0SKqe6/QFm773wKTBtlDec+zkitLw+wEWDAMx2xMbZW4RLf3ukTMjsLJjMaOQzX54V5ItB2O46kXTydQcpwffXQtp1+TrWg/nwkFWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAHWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAMWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAPWDAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABrY2VydGlmaWNhdGVZAn4wggJ6MIICAaADAgECAhABlPwqQ3vJMwAAAABnvL3OMAoGCCqGSM49BAMDMIGOMQswCQYDVQQGEwJVUzETMBEGA1UECAwKV2FzaGluZ3RvbjEQMA4GA1UEBwwHU2VhdHRsZTEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxOTA3BgNVBAMMMGktMDJhMjMxMDNhZDU4ZjRhNzgudXMtd2VzdC0yLmF3cy5uaXRyby1lbmNsYXZlczAeFw0yNTAyMjQxODQzMjNaFw0yNTAyMjQyMTQzMjZaMIGTMQswCQYDVQQGEwJVUzETMBEGA1UECAwKV2FzaGluZ3RvbjEQMA4GA1UEBwwHU2VhdHRsZTEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxPjA8BgNVBAMMNWktMDJhMjMxMDNhZDU4ZjRhNzgtZW5jMDE5NGZjMmE0MzdiYzkzMy51cy13ZXN0LTIuYXdzMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEkHUTUMrn0Cx/SjmnGQQx+uElzJN5GgoZtk8a4TO/7GM4n2VD/pytpd3fA+yVNVeqySdUGQ7VLr4zqmrAaDbvigk1kS1E9udtb8FIPqGTKDH9MrRoJDQkkPJv+Jcn8Glcox0wGzAMBgNVHRMBAf8EAjAAMAsGA1UdDwQEAwIGwDAKBggqhkjOPQQDAwNnADBkAjB0YPmYUhmaIztSxIJPFrL6PhaMf1GwbdcyObxhl0eJyOEG39S4t8QYdcROaWmD9lMCMF7fwje+yBIOxZhFg2Vf3o40GSdDFzOsSWVs6Id+l9Hm+RW7/XDs8mvwV4ura+SsM2hjYWJ1bmRsZYRZAhUwggIRMIIBlqADAgECAhEA+TF1aBuQr+EdRsy05Of4VjAKBggqhkjOPQQDAzBJMQswCQYDVQQGEwJVUzEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxGzAZBgNVBAMMEmF3cy5uaXRyby1lbmNsYXZlczAeFw0xOTEwMjgxMzI4MDVaFw00OTEwMjgxNDI4MDVaMEkxCzAJBgNVBAYTAlVTMQ8wDQYDVQQKDAZBbWF6b24xDDAKBgNVBAsMA0FXUzEbMBkGA1UEAwwSYXdzLm5pdHJvLWVuY2xhdmVzMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE/AJU66YIwfNocOKa2pC+RjgyknNuiUv/9nLZiURLUFHlNKSx9tvjwLxYGjK3sXYHDt4S1po/6iEbZudSz33R3QlfbxNw9BcIQ9ncEAEh5M9jASgJZkSHyXlihDBNxT/0o0IwQDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSQJbUN2QVH55bDlvpync+Zqd9LljAOBgNVHQ8BAf8EBAMCAYYwCgYIKoZIzj0EAwMDaQAwZgIxAKN/L5Ghyb1e57hifBaY0lUDjh8DQ/lbY6lijD05gJVFoR68vy47Vdiu7nG0w9at8wIxAKLzmxYFsnAopd1LoGm1AW5ltPvej+AGHWpTGX+c2vXZQ7xh/CvrA8tv7o0jAvPf9lkCwjCCAr4wggJFoAMCAQICEQCLtgI6xJuqK6j3iacjKfIPMAoGCCqGSM49BAMDMEkxCzAJBgNVBAYTAlVTMQ8wDQYDVQQKDAZBbWF6b24xDDAKBgNVBAsMA0FXUzEbMBkGA1UEAwwSYXdzLm5pdHJvLWVuY2xhdmVzMB4XDTI1MDIyMDIwMjMyMFoXDTI1MDMxMjIxMjMyMFowZDELMAkGA1UEBhMCVVMxDzANBgNVBAoMBkFtYXpvbjEMMAoGA1UECwwDQVdTMTYwNAYDVQQDDC00OGE2ZDFhZTcyZTlhMDczLnVzLXdlc3QtMi5hd3Mubml0cm8tZW5jbGF2ZXMwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAAR6+HeRSGgrxnOwkoxv8OKj2b3GbWqClLl6lLx4X7aJ6iPBJgX3WIkKt7o6qdH3tFaxhPO/TypAgEV/EjC5wpF5uDjSOUkOmCl3Nrj5zYbqFfspB/jmZuI5i/0f59pwrVajgdUwgdIwEgYDVR0TAQH/BAgwBgEB/wIBAjAfBgNVHSMEGDAWgBSQJbUN2QVH55bDlvpync+Zqd9LljAdBgNVHQ4EFgQUPrxPlv/nkqBUfRRgY2wpcACKQuQwDgYDVR0PAQH/BAQDAgGGMGwGA1UdHwRlMGMwYaBfoF2GW2h0dHA6Ly9hd3Mtbml0cm8tZW5jbGF2ZXMtY3JsLnMzLmFtYXpvbmF3cy5jb20vY3JsL2FiNDk2MGNjLTdkNjMtNDJiZC05ZTlmLTU5MzM4Y2I2N2Y4NC5jcmwwCgYIKoZIzj0EAwMDZwAwZAIwDjC93DcOXYeJWV8GHAnsN5x/DXG87IwbaD6flHCdQlAqlrrJDi16bI8xUcVDTGp3AjAq+93d4I45BxfuD4mon8CdcTujlqkcRHqd9qSmANBSv/UX6YUKbGmm9g7Myzooi1FZAxgwggMUMIICm6ADAgECAhEAw0RDGqXRvPYiiI+KlwsUzzAKBggqhkjOPQQDAzBkMQswCQYDVQQGEwJVUzEPMA0GA1UECgwGQW1hem9uMQwwCgYDVQQLDANBV1MxNjA0BgNVBAMMLTQ4YTZkMWFlNzJlOWEwNzMudXMtd2VzdC0yLmF3cy5uaXRyby1lbmNsYXZlczAeFw0yNTAyMjQwMDU1MDJaFw0yNTAzMDIwMTU1MDJaMIGJMTwwOgYDVQQDDDMzNmJiNzhkYzU4YjBhNWY1LnpvbmFsLnVzLXdlc3QtMi5hd3Mubml0cm8tZW5jbGF2ZXMxDDAKBgNVBAsMA0FXUzEPMA0GA1UECgwGQW1hem9uMQswCQYDVQQGEwJVUzELMAkGA1UECAwCV0ExEDAOBgNVBAcMB1NlYXR0bGUwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAASp9EzU6mADkRGK3cIMmaA9o6or6v0yWctk2nAVHgpEqgsCKTo8U/Q7f7UYmWQGOZMGTQ/Yqk246MxLdXUS0fLr+gQASRpy+LK48r7rwEPrWt8yqwliWIp1UNtFeXUDcEujgeowgecwEgYDVR0TAQH/BAgwBgEB/wIBATAfBgNVHSMEGDAWgBQ+vE+W/+eSoFR9FGBjbClwAIpC5DAdBgNVHQ4EFgQUN5HX6Ri+Ac+1PCknMiZfuEXdoVUwDgYDVR0PAQH/BAQDAgGGMIGABgNVHR8EeTB3MHWgc6Bxhm9odHRwOi8vY3JsLXVzLXdlc3QtMi1hd3Mtbml0cm8tZW5jbGF2ZXMuczMudXMtd2VzdC0yLmFtYXpvbmF3cy5jb20vY3JsLzM4MjgzMjFiLTBkNDgtNGZjMy04MTUyLWEzZjQ1YTNiMzBlZi5jcmwwCgYIKoZIzj0EAwMDZwAwZAIwGrRzhyqMrcFOge6buCdSIzfdrIZ+STLumfoC4o5FCpJ4tnMogDTVFHD5lKJ+dgSeAjBJrzaBs/sKWEbzdSjORdlRXui4iGeg2Wqnv0dDebBe5UFbdAXU+Zwo1DQar9wrO2pZAsIwggK+MIICRKADAgECAhQjyijWQmLyfnH8IKM4PjGeCXQ4GjAKBggqhkjOPQQDAzCBiTE8MDoGA1UEAwwzMzZiYjc4ZGM1OGIwYTVmNS56b25hbC51cy13ZXN0LTIuYXdzLm5pdHJvLWVuY2xhdmVzMQwwCgYDVQQLDANBV1MxDzANBgNVBAoMBkFtYXpvbjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAldBMRAwDgYDVQQHDAdTZWF0dGxlMB4XDTI1MDIyNDA5MjMzM1oXDTI1MDIyNTA5MjMzM1owgY4xCzAJBgNVBAYTAlVTMRMwEQYDVQQIDApXYXNoaW5ndG9uMRAwDgYDVQQHDAdTZWF0dGxlMQ8wDQYDVQQKDAZBbWF6b24xDDAKBgNVBAsMA0FXUzE5MDcGA1UEAwwwaS0wMmEyMzEwM2FkNThmNGE3OC51cy13ZXN0LTIuYXdzLm5pdHJvLWVuY2xhdmVzMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE3I1rZjXOUyVx4jHQIrLbj7NNeVLtZTvUPumM0G1zLfD9TFEDofUY4sfRTG+tqeiJe8G8pZ3djhMWmD+kNvir3vlV9efystofcIIFQB7+VuSBn8XHPHOpnqpIDr4KPUh8o2YwZDASBgNVHRMBAf8ECDAGAQH/AgEAMA4GA1UdDwEB/wQEAwICBDAdBgNVHQ4EFgQUVypfioMXTXFrZfqZD/cvaC4+xuUwHwYDVR0jBBgwFoAUN5HX6Ri+Ac+1PCknMiZfuEXdoVUwCgYIKoZIzj0EAwMDaAAwZQIwQU2Q4weVl5wU/sT/8w/eM53UI9jUp0STqHKFPMGrSV3ibv0emBZX/oGdmtA0ly5QAjEA+yb1Q0wnZoyUeP+b5i6LibCCNHsP1+Espx3cEB/GBdMixCy5GrjzY3Rtjm4CzrHdanB1YmxpY19rZXn2aXVzZXJfZGF0YVgg5/2xJ0M8FiX2xSlhqzwi2myru1mtlAcqvtyTSkdlzUdlbm9uY2X2WGD7GzI9ZqwZIuTf/L92bcSRUW4TWa3KLRpVNqLMiHW3eXer1Rv6m2hiCc0GBhyEmf9bfa5GN1O1lvVQhxLT5wEIhEzcWAA21XPwfWvs8AEjv29Bp/2aIu3Lkn8JVRuoZa0="
      ]
    },
    "public_key": {
      "curve_type": "p256k1",
      "data": "BFNSuXpjDrvnZOnIy6hPjJGs40oMCXEiKsLt87CIEmrX/rjmpA7E8+Y5Lq7tQeoB8HRwY1BNnStcWYBisdXQHEE="
    }
  },
  "function_calls": [
    {
      "transitive_attestation": "WyJXeUphYlZVMFRYcGpNbHBVVW1sT1JHY3dUVzFXYVU5VVZtdFBWMGw0VFVSV2ExcFhXbXRPUkdjeVRWUkZlazVVVVhoYVZHc3pUbFJDYlUxVVkzcE5SMUp0VG5wb2FrNXFTVFZOZWswMFRWUkplRmw2VlhsT1Yxa3dUbXBqZUU1SFdtbFpiVlV5VG5wS2JVNUVSVE5hYWtreFRtcEZlRTFFVFROYVJFVXdUbFJLYWs5RWJHdE9hbHBzV2xSUmVVOUhUVFJQUkZrMVRrUkJNRnBxWkcxTmJVNXNXVlJuTVUxWFdtcFphbEU5SWl3aVlVZFdjMkpIT1VaamJrcDJZMmM5UFNJc0lsbFVXVFZhYW1ONldUSk9hRTFxVG1oUFYwWnFUbGROTkZscVZUSk9NbEpxVFZSbk1WbFVZekZPYlZVMVRqSk5OVTlFU1hoT2FsSnRXbFJKTVU5RVZUVmFWRUpyVFZkU2FsbDZSVEJPZWxacVQwUkNhRTVxUlRGWmFrbDRUV3BPYUZwcVJtMU9WMWsxVGtkTmVFMVhWWHBhVkdzd1RVUkthazB5Um1wT1ZGVTBXbXBWZDAxRVJUVlBWMUUxVGxkSk1scEVUbXhOZWtGNFRucFZORTVVWnpKTmFtZDRXa2RPYTAxcVdUMGlMQ0psZVVwVVpGZE9hbHBZVG5wSmFuQnRXVmQ0ZWxwVGQybFNXRXA1WWpOSmFVOXBTbXhsU0VKc1dUTlNiRnBEUW14amJrcDJZMmxKYzBsck9URmtTRUl4WkVOSk5tSnVWbk5pU0RBOUlpd2lUMVJOTTA1VVVUQk9NazVyVGxSTmQwNHlTbTFPZWxFelRUSkpORTFxUVhkYWFrRjZUMWRKTWsxSFJYcFpiVlV3VDFSRmVVOUVTbTFQUkZWNVdrZFpOVnBxVVhsWk1sVjZUVmRGTkZsVVVYcGFhbHB0VDBkVk5VMVVXbXBPUjFrMFRXcFpNRnBVWkd0TmFrMTZXVmRTYTA1RVp6Tk9SRnBvVGtSQmVFNXFXbXhhVjAweFQwUm9hVnBVYUdsT01razFXV3BGTWxsVVZteFphbGsxVDBkUk1GbDZUbWxOUkZwc1RVUkJQU0pkIiwiWEFTOXJCZjVzemxoUUo3TTQxbkwxTWdMREUzeHYxZmw5OVlocWZ1dGtlZ3ZhRlVtOTZRbTJ4WnBzdkVuVkI5cExJWDd4R1h4OUIzbmgrRHZqS3JqRGdBPSJd",
      "claims": {
        "hash_of_code": "fe8376e4b4842eb95d9b105defd486113541e9750f1730df78c629338121c525f46714fbbe672f417f25611037d1452c89d66ee428c8869404f7f2cea851fcb4",
        "function": "helloError",
        "hash_of_input": "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26",
        "output": "eyJTdWNjZXNzIjpmYWxzZSwiRXJyb3IiOiJleHBlY3RlZCBlcnJvciIsIk91dHB1dCI6bnVsbH0=",
        "hash_of_secrets": "9375447cd5307bf7473b8200f039b60a3be491282f852df9f42ce31a8a43f6f8e916c4f8264e7d233add48746a40166eec588be8b7b9b16a5eb698d4c3b06e00"
      },
      "logs": "UmV0dXJuaW5nIGFuIGV4cGVjdGVkIGVycm9y"
    }
  ]
}
```

There you will see the `enclave_attested_application_public_key` that contains
the `enclave_attestation` over the Blocky AS server public key.
You will also see the `function_calls` section that contains the `transitive_attestation` over the function call.
The `bky-as` CLI verifies the `enclave_attestation`, extracts the Blocky AS
server public key and uses it to verify the `transitive_attestation` to extract
the `claims`.
You can learn more about this process in the
[Attestations in the Blocky Attestation Service](https://blocky-docs.redocly.app/attestation-service/concepts#attestations-in-the-blocky-attestation-service) section in our documentation.

The `claims` section contains attested information about the execution of
the function.
You can see:

- `hash_of_code`: The hash of the WASM code executed by the Blocky AS server
- `function`: The name of the function executed by the Blocky AS server
- `hash_of_input`: The hash of the input data used by the function. In this
  example this is the hash of the empty string, since we didn't specify any
  input.
- `hash_of_secrets`: The hash of the secrets used by the function. In this
  example this is the hash of the empty string, since we didn't specify any
  secrets.
- `output`: The output of the function encoded in base64.

Finally, the `logs` field contains the logs generated by the function execution
and encoded in base64.

If you look at the `run` target in the [Makefile](./Makefile) again, you will
see that we use `jq` to extract the `output` and the `logs` fields from
`tmp/out.json` and decode them from base64.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. Check out other examples in this repository, to learn what
else you can do with Blocky AS.


