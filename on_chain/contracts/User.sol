// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {TAParserLib} from "../lib/TAParserLib.sol";

contract User {
    event AttestedFunctionCallOutput(string output);

    address private enclAttAppPubKeyAddress;

    function setEnclaveAttestedAppPubKey(
        bytes calldata enclAttAppPubKey
    )
    public
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
