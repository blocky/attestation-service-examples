// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

import {JsmnSolLib} from "../lib/JsmnSolLib.sol";
import {TAParser} from "./TAParser.sol";
import {console} from "hardhat/console.sol";

contract User is TAParser {
    event AttestedFunctionCallOutput(string result, string error);

    address public _verifierAddress;
    address public taSigningKeyAddress;

    // todo: make this an owner function
    function setTASigningKeyAddress(
        bytes calldata taSigningKey
    )
    public
    {
        taSigningKeyAddress = publicKeyToAddress(taSigningKey);
    }

    function verifyAttestedAPICallClaims(
        string calldata taData
    )
    private view returns (TAParser.FnCallClaims memory)
    {
        TAParser.FnCallClaims memory claims = parseTA(
            taData,
            taSigningKeyAddress
        );

        return claims;
    }

    function parseAPICallClaims(
        TAParser.FnCallClaims memory claims
    ) public
    {
        string memory out = base64d(claims.Output);

        JsmnSolLib.Token[] memory tokens;
        uint number;
        uint success;
        (success, tokens, number) = JsmnSolLib.parse(out, 50);

        uint isErrIdx = 2;
        string memory isErr = JsmnSolLib.getBytes(out, tokens[isErrIdx].start, tokens[isErrIdx].end);

        uint avgIdx = 6;
        string memory avg = JsmnSolLib.getBytes(out, tokens[avgIdx].start, tokens[avgIdx].end);

        console.log("Did it Error:", isErr);
        console.log("Average:", avg);

        emit AttestedFunctionCallOutput(isErr, avg);
    }

    function processAttestedAPICallClaims(
        string calldata taData
    ) public {
        console.log("\n> Processing attested API call claims");

        TAParser.FnCallClaims memory claims = verifyAttestedAPICallClaims(taData);

        parseAPICallClaims(claims);

        console.log("Processed attested API call claims");
    }
}

