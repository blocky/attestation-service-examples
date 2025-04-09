//  SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {JsmnSolLib} from "./JsmnSolLib.sol";
import {Base64} from "base64-sol/base64.sol";
import {BytesLib} from "solidity-bytes-utils/contracts/BytesLib.sol";

library TAParserLib {

    struct TA {
        string Data;
        string Sig;
    }

    struct FnCallClaims {
        string HashOfCode;
        string Function;
        string HashOfInput;
        string HashOfSecrets;
        string Output;
    }

    function publicKeyToAddress(
        bytes calldata publicKey
    )
        internal pure returns (address)
    {
        // strip out the public key prefix byte
        bytes memory strippedPublicKey = new bytes(publicKey.length - 1);
        for (uint i = 0; i < strippedPublicKey.length; i++) {
            strippedPublicKey[i] = publicKey[i + 1];
        }

        return address(uint160(uint256(keccak256(strippedPublicKey))));
    }

    function verifyTransitivelyAttestedFnCall(
        address applicationPublicKey,
        string calldata transitiveAttestation
    )
        internal pure returns (FnCallClaims memory)
    {
        string memory verifiedTAData = verifyTA(
            applicationPublicKey,
            transitiveAttestation
        );
        TAParserLib.FnCallClaims memory verifiedClaims = decodeFnCallClaims(
            verifiedTAData
        );
        return verifiedClaims;
    }

    function verifyTA(
        address publicKeyAddress,
        string calldata transitiveAttestation
    )
        private pure returns (string memory)
    {
        TA memory ta = decodeTA(transitiveAttestation);

        bytes memory sigAsBytes = Base64.decode(ta.Sig);
        bytes32 r = BytesLib.toBytes32(sigAsBytes, 0);
        bytes32 s = BytesLib.toBytes32(sigAsBytes, 32);
        uint8 v = 27 + uint8(sigAsBytes[64]);

        bytes memory dataAsBytes = Base64.decode(ta.Data);
        bytes32 dataHash = keccak256(dataAsBytes);
        address recovered = ecrecover(dataHash, v, r, s);

        require(publicKeyAddress == recovered, "Could not verify signature");

        return ta.Data;
    }

    function decodeTA(
        string calldata taData
    )
        private pure returns (TA memory)
    {
        TA memory ta;

        JsmnSolLib.Token[] memory tokens;
        uint number;
        uint success;
        (success, tokens, number) = JsmnSolLib.parse(taData, 3);

        // The flat rep of a TA is an array of 2 byte arrays.
        // As such, we expect to have 3 tokens.
        // Item 0 is the array of byte arrays with 1 and 2 as the content.
        require(number == 3, "Wrong number of tokens");
        ta.Data = JsmnSolLib.getBytes(taData, tokens[1].start, tokens[1].end);
        ta.Sig = JsmnSolLib.getBytes(taData, tokens[2].start, tokens[2].end);

        return ta;
    }

    function decodeFnCallClaims(
        string memory data
    )
        private pure returns (FnCallClaims memory)
    {
        FnCallClaims memory claims;

        string memory b64 = base64d(data);

        JsmnSolLib.Token[] memory tokens;
        uint number;
        uint success;
        (success, tokens, number) = JsmnSolLib.parse(b64, 6);

        // The flat rep of fn call claims is an array of 5 byte arrays.
        // As such, we expect to have 6 tokens.
        // Item 0 is the array of byte arrays with 1..5 as the content.
        require(number == 6, "Wrong number of tokens");
        claims.HashOfCode = base64d(JsmnSolLib.getBytes(b64, tokens[1].start, tokens[1].end));
        claims.Function = base64d(JsmnSolLib.getBytes(b64, tokens[2].start, tokens[2].end));
        claims.HashOfInput = base64d(JsmnSolLib.getBytes(b64, tokens[3].start, tokens[3].end));
        claims.Output = base64d(JsmnSolLib.getBytes(b64, tokens[4].start, tokens[4].end));
        claims.HashOfSecrets = base64d(JsmnSolLib.getBytes(b64, tokens[5].start, tokens[5].end));

        return claims;
    }

    function base64d(
        string memory base64Input
    )
        private pure returns (string memory)
    {
        bytes memory decodedBytes = Base64.decode(base64Input);
        return string(decodedBytes);
    }
}
