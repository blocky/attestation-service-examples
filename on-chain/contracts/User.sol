// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/access/Ownable.sol";
import "hardhat/console.sol";
import {JsmnSolLib} from "../lib/JsmnSolLib.sol";
import {TAParser} from "./TAParser.sol";
import {console} from "hardhat/console.sol";

contract User is Ownable, TAParser {
    event AttestedFunctionCallOutput(string Success, string TWAP);

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

        uint successIdx = 4;
        string memory resultSuccess = JsmnSolLib.getBytes(out, tokens[successIdx].start, tokens[successIdx].end);

        uint twapIdx = 8;
        string memory resultTWAP = JsmnSolLib.getBytes(out, tokens[twapIdx].start, tokens[twapIdx].end);

        console.log("Success:", resultSuccess);
        console.log("TWAP:", resultTWAP);

        emit AttestedFunctionCallOutput(resultSuccess, resultTWAP);
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

