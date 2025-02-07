# Iteration

This demo shows how we feed attested outputs as input into the attestation
service.  One great use for that features is that it allows us to carry state
from one invocation to the next.  For example, here, we iteratively call the
Coingecko API to accumulate a collection of samples such that we can compute
the average of those samples.

To get started, put a copy of `bky-as` in the current directory.  You can use
`nix` to set up the correct shell (but note that it sets up an alias that is
pretty specific to dave's prototyping environment).

```bash
nix-shell
```
Build our wasm function.  Note that you may see an error 

```bash
jq: error: Could not open file tmp/prev.json: No such file or directory
```

No worries, that just means we have no previous iteration (and because this is
a quick demo, I didn't try to fix that error message.)


```bash
make build
```

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


