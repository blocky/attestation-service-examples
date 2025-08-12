// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {TAParserLib} from "../lib/TAParserLib.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract User is Ownable {
    event AttestedFunctionCallOutput(string output);

    address private enclAttAppPubKeyAddress;

    constructor() Ownable(msg.sender) {
    }

    function setEnclaveAttestedAppPubKey(
        bytes calldata enclAttAppPubKey
    )
    public onlyOwner
    {
        enclAttAppPubKeyAddress = TAParserLib.publicKeyToAddress(
            enclAttAppPubKey
        );
    }

    function processTransitiveAttestedFunctionCall(
        bytes calldata transitiveAttestation
    )
    public
    {
        TAParserLib.FnCallClaims memory claims;

        claims = TAParserLib.verifyTransitiveAttestedFnCall(
            enclAttAppPubKeyAddress,
            transitiveAttestation
        );

        emit AttestedFunctionCallOutput(string(claims.Output));
    }
}
