# On Chain

This demo shows how we can get transitive attested data on chain. In this demo,
bring a transitive attested function call from the [twap demo](../attest_fn_call)
on to a chain in a development environment.

## Setup

Set up the project dependencies:

```bash
npm install
```

## Run

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
