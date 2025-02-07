# On Chain

This demo shows how we can get transitive attested data on chain. In this demo,
bring a transitive attested function call from the [iterate demo](../iterate/)
on to a chain in a development environment.

## Setup

You should not do this again, but here are the steps documenting getting
the project setup.

```bash
npm install npm@10.9.0
npm init
npm install --save-dev hardhat
npx hardhat init
```

The last step may hang, you can cancel and reissue the command

```bash
npm install --save-dev "@nomicfoundation/hardhat-toolbox@^5.0.0" --dd
```

Install other dependencies


```bash
npm install --save-dev -dd solidity-bytes-utils
```

You are ready to start going!!

## Bring results from iterate on chain

Copy an output from iterate

```bash
mkdir input
cp ../iterate/tmp/prev.json input/prev.json
```

And run the test:

```bash
make test-local
```
