//  SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import {JsmnSolLib} from "./JsmnSolLib.sol";
import {Base64} from "base64-sol/base64.sol";
import {BytesLib} from "solidity-bytes-utils/contracts/BytesLib.sol";
import {console} from "hardhat/console.sol";

library TAParserLib {

    struct TA {
        bytes Data;
        bytes Sig;
    }

    struct FnCallClaims {
        bytes HashOfCode;
        bytes Function;
        bytes HashOfInput;
        bytes HashOfSecrets;
        bytes Output;
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
        bytes calldata transitiveAttestation
    )
        internal pure returns (FnCallClaims memory)
    {
        bytes memory verifiedTAData = verifyTA(
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
        bytes calldata transitiveAttestation
    )
        private pure returns (bytes memory)
    {
        TA memory ta = decodeTA(transitiveAttestation);

        // bytes memory sigAsBytes = Base64.decode(ta.Sig);
        bytes memory sigAsBytes = ta.Sig;
        bytes32 r = BytesLib.toBytes32(sigAsBytes, 0);
        bytes32 s = BytesLib.toBytes32(sigAsBytes, 32);
        uint8 v = 27 + uint8(sigAsBytes[64]);

        // bytes memory dataAsBytes = Base64.decode(ta.Data);
        bytes memory dataAsBytes = ta.Data;
        bytes32 dataHash = keccak256(dataAsBytes);
        address recovered = ecrecover(dataHash, v, r, s);

        require(publicKeyAddress == recovered, "Could not verify signature");

        return ta.Data;
    }


     function decodeInput(bytes calldata input) external pure returns (bytes memory, bytes memory) {
            bytes[] memory decoded = abi.decode(input, (bytes[]));
            return (decoded[0], decoded[1]);
        }


    function decodeTA(
        bytes calldata taData
    )
        private pure returns (TA memory)
    {
        TA memory ta;

        bytes[] memory decodedTA = abi.decode(taData, (bytes[]));
        require(decodedTA.length == 2, "Expected 2 elements");

        ta.Data = decodedTA[0];
        ta.Sig = decodedTA[1];

        return ta;
    }

    function decodeFnCallClaims(
        bytes memory data
    )
        private pure returns (FnCallClaims memory)
    {
        FnCallClaims memory claims;

        bytes[] memory decodedData = abi.decode(data, (bytes[]));
        require(decodedData.length == 5, "Expected 5 elements");

        claims.HashOfCode = decodedData[0];
        claims.Function = decodedData[1];
        claims.HashOfInput = decodedData[2];
        claims.Output = decodedData[3];
        claims.HashOfSecrets = decodedData[4];

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
