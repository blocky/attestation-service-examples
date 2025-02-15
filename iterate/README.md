# Iteration

This demo shows how we feed attested outputs as input into the attestation
service. One great use for that features is that it allows us to carry state
from one invocation to the next. For example, here, we iteratively call the
Coingecko API to accumulate a collection of samples such that we can compute
the average of those samples.

To get started, install `bky-as` in the current directory by following the
instructions in the
[Blocky Attestation Service setup documentation](https://blocky-docs.redocly.app/attestation-service/setup).

Next, use `nix` to set up the correct environment for building WASM binaries 
with TinyGo. Note that this process will take a while the first time you run it.

```bash
nix-shell
```

Build our WASM function.

```bash
make build
```

Note that you may see an error

```bash
jq: error: Could not open file tmp/prev.json: No such file or directory
```

No worries, that just means we have no previous iteration (and because this is
a quick demo, we didn't try to fix that error message.)

Init the iteration

```bash
make init
```

Ran an iteration

```bash
make run
```

And another...

```bash
make run
```

When you are all done, clean up with:

```bash
make clean
```


