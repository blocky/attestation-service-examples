// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {JsmnSolLib} from "../lib/JsmnSolLib.sol";
import {TAParserLib} from "../lib/TAParserLib.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {console} from "hardhat/console.sol";

contract User is Ownable {
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
        taSigningKeyAddress = TAParserLib.publicKeyToAddress(taSigningKey);
    }

    function verifyAttestedFnCallClaims(
        string calldata taData
    )
    private view returns (TAParserLib.FnCallClaims memory)
    {
        TAParserLib.FnCallClaims memory claims = TAParserLib.parseTA(
            taData,
            taSigningKeyAddress
        );

        return claims;
    }

    function parseFnCallClaims(
        TAParserLib.FnCallClaims memory claims
    ) public
    {
        JsmnSolLib.Token[] memory tokens;
        uint number;
        uint success;
        (success, tokens, number) = JsmnSolLib.parse(claims.Output, 50);

        uint successIdx = 2;
        bool resultSuccess = JsmnSolLib.parseBool(
            JsmnSolLib.getBytes(
                claims.Output,
                tokens[successIdx].start,
                tokens[successIdx].end
            )
        );

        uint errorIdx = 4;
        string memory resultError = JsmnSolLib.getBytes(
            claims.Output,
            tokens[errorIdx].start,
            tokens[errorIdx].end
        );

        require(resultSuccess, resultError);

        uint twapIdx = 6;
        string memory resultTWAP = JsmnSolLib.getBytes(
            claims.Output,
            tokens[twapIdx].start,
            tokens[twapIdx].end
        );

        emit AttestedFunctionCallOutput(resultTWAP);
    }

    function processAttestedFnCallClaims(
        string calldata taData
    ) public {
        console.log("\n> Processing attested function call claims");

        TAParserLib.FnCallClaims memory claims = TAParserLib.parseTA(
            taData,
            taSigningKeyAddress
        );

        parseFnCallClaims(claims);

        console.log("Processed attested function call claims");
    }
}

