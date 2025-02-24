# Hello World Example

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a simple function calls.

You'll learn how to:
- Create a function that returns a "Hello, World!" message
- Create a simple function that writes a message to a log and returns a
user-defined error
- Invoke functions in the Blocky AS using its `bky-as` CLI
- Process Blocky AS attestations

## Setup

- Install the Blocky AS CLI by following the
[setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
in the Blocky AS documentation. 
- Make sure you also have
[Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
your system.
- [Get a key for the CoinGecko API](https://docs.coingecko.com/reference/setting-up-your-api-key) 
and set it in `fn-call.json` in the `api_key` field. 

## Quick Start

To run this example, call the following command:

```bash
make FUNCTION=hello-world
```

You will see the following output:

```json
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

