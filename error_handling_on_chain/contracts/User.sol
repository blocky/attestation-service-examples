// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {JsmnSolLib} from "../lib/JsmnSolLib.sol";
import {TAParserLib} from "../lib/TAParserLib.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {console} from "hardhat/console.sol";

contract User is Ownable {
    event ResultValue(string output);

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

        parseResult(claims.Output);
    }

    function parseResult(
        string memory resultString
    ) public
    {
        JsmnSolLib.Token[] memory tokens;
        uint number;
        uint success;
        (success, tokens, number) = JsmnSolLib.parse(resultString, 50);

        bool resultSuccess;
        string memory resultError;
        string memory valueString;

        for (uint i = 0; i < number; i++) {
            if (tokens[i].jsmnType == JsmnSolLib.JsmnType.STRING) {
                string memory key = JsmnSolLib.getBytes(
                    resultString,
                    tokens[i].start,
                    tokens[i].end
                );

                if (keccak256(bytes(key)) == keccak256("Success")) {
                    resultSuccess = JsmnSolLib.parseBool(
                        JsmnSolLib.getBytes(
                            resultString,
                            tokens[i + 1].start,
                            tokens[i + 1].end
                        )
                    );
                } else if (keccak256(bytes(key)) == keccak256("Error")) {
                    resultError = JsmnSolLib.getBytes(
                        resultString,
                        tokens[i + 1].start,
                        tokens[i + 1].end
                    );
                } else if (keccak256(bytes(key)) == keccak256("Value")) {
                    valueString = JsmnSolLib.getBytes(
                        resultString,
                        tokens[i + 1].start,
                        tokens[i + 1].end
                    );
                }
            }
        }

        console.log("\tSuccess: %s", resultSuccess);
        console.log("\tError: %s", resultError);
        console.log("\tValue: %s", valueString);

        require(resultSuccess, resultError);
        emit ResultValue(valueString);
    }
}

