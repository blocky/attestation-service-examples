// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {TAParserLib} from "../lib/TAParserLib.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {console} from "hardhat/console.sol";

contract User is Ownable {
    event AttestedFunctionCallOutput(string output);

    address private taSigningKeyAddress;

    constructor() Ownable(msg.sender) {
    }

    function setTASigningKeyAddress(
        bytes calldata taSigningKey
    )
    public onlyOwner
    {
        taSigningKeyAddress = TAParserLib.publicKeyToAddress(taSigningKey);
    }

    function verifyAttestedFnCallClaims(
        string calldata taData
    )
    public
    {
        TAParserLib.FnCallClaims memory claims = TAParserLib.parseTA(
            taData,
            taSigningKeyAddress
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

