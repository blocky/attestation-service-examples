# On Chain

This demo shows how we can get transitive attested data on chain. In this demo,
bring a transitive attested function call from the [twap demo](../twap)
on to a chain in a development environment.

## Setup

Use `nix` to set up a shell environment with the correct dependencies:

```bash
nix-shell
```

Continue with the following commands in the nix shell.

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

You can see the details of this test in [test/User.ts](test/User.ts).

(Optionally) If you ran the [twap demo](../twap) and want to use its
latest transitive attested function call, bring it into this project by
running:

```bash
make copy-from-iterate
```
