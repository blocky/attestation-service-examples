# Bringing An Attestation Of Time-Weighed Average Price On Chain

This example shows you how to bring an attestation of time-weighted average price (TWAP)
on chain. It builds on the 
[Attesting A Time-Weighted Average Price](../attest_fn_call) 
example, which shows you how to obtain an attestation of a TWAP, and on the
[Bringing A Blocky AS Function Call Attestation On Chain](../../on_chain)
example, which shows you how to bring an attestation of a Blocky AS function
call on chain.

In this example, you'll learn how to:

- Bring an attestation of a TWAP on chain 

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

which will test verifying an attested TWAP in
[`contracts/User.sol`](contracts/User.sol)
within a local test environment and give your output like:

```
  Local Tests
    ✔ Verify attested TWAP in User contract (1080ms)
```

## Walkthrough

We cover the basics of how to verify an attested function call in a smart 
contract in the
[Bringing A Blocky AS Function Call Attestation On Chain](../../on_chain)
example. 

In this example, we go a step further and show you how to parse out the 
TWAP value from the `Result` struct attested in the output of the `twap` function
call in the [Attesting A Time-Weighted Average Price](../attest_fn_call) example.
If you would like more background on how we use the `Result` struct, please 
review the [Error Handling - Attested Function Calls](../../error_handling_attest_fn_call) example.

In [`contracts/User.sol`](contracts/User.sol), we define a `parseTWAP` function:

```solidity
    function parseTWAP(
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

    require(resultSuccess, resultError);
    emit TWAP(valueString);
}
```

which takes in the verified `TAParserLib.FnCallClaims` `claims` and uses
`JsmnSolLib` to parse the JSON serialized `Result` struct contains in
`claims.Output`. We use positional arguments to parse out `resultSuccess` from
`Result.Success` and `resultError` from `Result.Error` and if the
`resultSuccess` is `false`, we revert the transaction with the `resultError`
message. If `resultSuccess` is `true`, we parse out the `resultTWAP` from
`Result.Value` and emit it as a `TWAP` event.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. Check out other examples in this repository, to learn what
else you can do with Blocky AS.
