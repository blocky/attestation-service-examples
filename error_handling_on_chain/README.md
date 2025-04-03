# Error Handling On Chain

This example shows you how to handle errors in Blocky AS function calls on chain.
It builds on the [Hello World On Chain](../hello_world_on_chain) and the
[Error Handling Attest Function Call](../error_handling_attest_fn_call)
examples, which show you how to verify and parse Blocky AS attestations on
chain and handle errors returned by attested function calls.

In this example, you'll learn how to:

- Write a smart contract to extract the `Result` from a function attestation
- Check `Result` for an error to determine whether to return the output or to
revoke the transaction

## Setup

Set up the project dependencies:

```bash
npm install
```

## Quick Start

To run this example, call:

```bash
make test-local
```

which will test extracting and verifying `Result` from an attested fn call in
[`contracts/User.sol`](contracts/User.sol) 
within a local test environment:

```
  Local Test
        Success: true
        Error: 
        Value: {"number":42}
    ✔ Verify TA and parse Result w/success (616ms)
        Success: false
        Error: expected error
        Value: null
    ✔ Verify TA and parse Result w/error (189ms)
```

## Walkthrough

We cover the basics of how to verify an attested function call in a smart
contract in the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../../hello_world_on_chain)
example.

In this example, we go a step further and show you how to extract the `Result`
struct from the `successFunc` and `errorFunc` function call attestations in the
[Error Handling](../error_handling_attest_fn_call) example.

In [`contracts/User.sol`](contracts/User.sol), we define a `parseResult` function:

```solidity
function parseResult(
    string memory resultString
) public
{
    JsmnSolLib.Token[] memory tokens;
    uint number;
    uint success;
    (success, tokens, number) = JsmnSolLib.parse(resultString, 50);

    bool resultSuccess;
    string memory resultError;
    string memory valueString;

    for (uint i = 0; i < number; i++) {
        if (tokens[i].jsmnType == JsmnSolLib.JsmnType.STRING) {
            string memory key = JsmnSolLib.getBytes(
                resultString,
                tokens[i].start,
                tokens[i].end
            );

            if (keccak256(bytes(key)) == keccak256("Success")) {
                resultSuccess = JsmnSolLib.parseBool(
                    JsmnSolLib.getBytes(
                        resultString,
                        tokens[i + 1].start,
                        tokens[i + 1].end
                    )
                );
            } else if (keccak256(bytes(key)) == keccak256("Error")) {
                resultError = JsmnSolLib.getBytes(
                    resultString,
                    tokens[i + 1].start,
                    tokens[i + 1].end
                );
            } else if (keccak256(bytes(key)) == keccak256("Value")) {
                valueString = JsmnSolLib.getBytes(
                    resultString,
                    tokens[i + 1].start,
                    tokens[i + 1].end
                );
            }
        }
    }

    console.log("\tSuccess: %s", resultSuccess);
    console.log("\tError: %s", resultError);
    console.log("\tValue: %s", valueString);

    require(resultSuccess, resultError);
    emit ResultValue(valueString);
}
```

which takes in the verified `TAParserLib.FnCallClaims` `claims.Output` and uses
`JsmnSolLib` to parse the JSON serialized `Result` struct. JSON serialization
does not guarantee struct field order, so we loop over the tokens to extract the
fields: `Result.Success`, `Result.Error`, and `Result.Value`. If `Result.Success`
is `false`, we revert the transaction with `Result.Error` as the error message.
If `Result.Success` is `true`, we emit `Result.Value` as an event, where
`Result.Value` itself may be a JSON serialized struct representing the output
of the function call.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For instance, write a function that extracts the`Output`
struct from the `ResultValue(valueString)` event. Check out other examples in
this repository, to learn what else you can do with Blocky AS.
