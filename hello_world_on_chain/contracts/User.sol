// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {TAParserLib} from "../lib/TAParserLib.sol";
import {console} from "hardhat/console.sol";

contract User {
    event AttestedFunctionCallOutput(string output);

    function labeledLog(
        string memory label,
        bytes memory data
    )
    public pure
    {
        console.log("\t%s: %s", label, string(data));
    }

    function processTransitivelyAttestedHelloWorldOutput(
        bytes calldata applicationPublicKey,
        bytes calldata transitiveAttestation
    )
    public
    {
        TAParserLib.FnCallClaims memory claims;

        address applicationPublicKeyAsAddress = TAParserLib.publicKeyToAddress(
            applicationPublicKey
        );

        claims = TAParserLib.verifyTransitivelyAttestedFnCall(
            applicationPublicKeyAsAddress,
            transitiveAttestation
        );

        console.log("Verified attest-fn-call claims:");
        labeledLog("Function", claims.Function);
        labeledLog("Hash of code",claims.HashOfCode);
        labeledLog("Hash of input", claims.HashOfInput);
        labeledLog("Hash of secrets", claims.HashOfSecrets);
        labeledLog("Output,", claims.Output);

        emit AttestedFunctionCallOutput(string(claims.Output));
    }
}
