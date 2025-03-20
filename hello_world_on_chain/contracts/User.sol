// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {TAParserLib} from "../lib/TAParserLib.sol";
import {console} from "hardhat/console.sol";

contract User {
    event AttestedFunctionCallOutput(string output);

        bytes calldata publicKey,
    function processTAHelloWorld(
        string calldata transitiveAttestation
    )
        public
    {
        TAParserLib.FnCallClaims memory claims = TAParserLib.verifyAttestedFnCall(
            publicKey,
            transitiveAttestation
        );

        console.log("Verified attest-fn-call claims:");
        console.log("\tFunction: %s", claims.Function);
        console.log("\tHash of code: %s", claims.HashOfCode);
        console.log("\tHash of input: %s", claims.HashOfInput);
        console.log("\tHash of secrets: %s", claims.HashOfSecrets);
        console.log("\tOutput: %s", claims.Output);

        emit AttestedFunctionCallOutput(claims.Output);
    }
}
