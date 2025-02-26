// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {JsmnSolLib} from "../lib/JsmnSolLib.sol";
import {TAParser} from "./TAParser.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {console} from "hardhat/console.sol";

contract User is Ownable, TAParser {
    event AttestedFunctionCallOutput(string output);

    address public _verifierAddress;
    address public taSigningKeyAddress;

    constructor() Ownable(msg.sender) {
    }

    function setTASigningKeyAddress(
        bytes calldata taSigningKey
    )
    public onlyOwner
    {
        taSigningKeyAddress = publicKeyToAddress(taSigningKey);
    }

    function verifyAttestedFnCallClaims(
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

    function parseFnCallClaims(
        TAParser.FnCallClaims memory claims
    ) public
    {
        string memory out = base64d(claims.Output);

        JsmnSolLib.Token[] memory tokens;
        uint number;
        uint success;
        (success, tokens, number) = JsmnSolLib.parse(out, 50);

        uint successIdx = 2;
        bool resultSuccess = JsmnSolLib.parseBool(
            JsmnSolLib.getBytes(
                out,
                tokens[successIdx].start,
                tokens[successIdx].end
            )
        );

        uint errorIdx = 4;
        string memory resultError = JsmnSolLib.getBytes(
            out,
            tokens[errorIdx].start,
            tokens[errorIdx].end
        );

        require(resultSuccess, resultError);

        uint outputIdx = 6;
        string memory resultOutput = JsmnSolLib.getBytes(
            out,
            tokens[outputIdx].start,
            tokens[outputIdx].end
        );

        emit AttestedFunctionCallOutput(resultOutput);
    }

    function processAttestedFnCallClaims(
        string calldata taData
    ) public {
        console.log("\n> Processing attested function call claims");

        TAParser.FnCallClaims memory claims = verifyAttestedFnCallClaims(taData);

        parseFnCallClaims(claims);

        console.log("Processed attested function call claims");
    }
}

