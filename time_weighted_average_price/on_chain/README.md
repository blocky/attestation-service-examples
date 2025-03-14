# Bringing An Attestation Of Time-Weighed Average Price On Chain

This example shows you how to bring an attestation of time-weighted average price (TWAP)
on chain. It builds on the 
[Attesting A Time-Weighted Average Price](../attest_fn_call) 
example, which shows you how to obtain an attestation of a TWAP, and on the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../../hello_world_on_chain)
example, which shows you how to bring an attestation of a Blocky AS function
call on chain.

In this example, you'll learn how to:

- Bring an attestation of a TWAP on chain 
- Revert a smart contract transaction if the 

## Setup

Set up the project dependencies:

```bash
npm install --dd
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

> Processing attested function call claims
Processed attested function call claims
    ✔ process TA (975ms)
```

## Walkthrough

We cover the details of how to verify an attested function call in a smart 
contract in the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../../hello_world_on_chain)
example.
This example works in 


Test bringing a transitive attested function call on chain into the
[User](contracts/User.sol) contract:

```bash
make test-local
```

You can see the details of this test in [test/user.ts](test/user.ts).

(Optionally) If you ran the [twap demo](../attest_fn_call) and want to use its
latest transitive attested function call, bring it into this project by
running:

```bash
make copy-from-twap
```
