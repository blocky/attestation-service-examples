// noinspection DuplicatedCode

import hre from "hardhat";
import {loadFixture} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import {expect} from "chai";
import {ethers} from "ethers";
import fs from "fs";
import path from "path";
import {User} from "../typechain-types";

type EVMLinkData = {
    publicKey: string;
    transitiveAttestation: string;
};

function loadEVMLinkData(jsonPath: string): EVMLinkData {
    try {
        const dir: string = path.resolve(__dirname, jsonPath);
        const file: string = fs.readFileSync(dir, "utf8");

        const data: any = JSON.parse(file);

        const k: any =
            data.enclave_attested_application_public_key.claims.public_key.data
        const pubKeyBytes: Uint8Array = ethers.decodeBase64(k)
        const publicKeyHex: string = Buffer.from(pubKeyBytes).toString('hex');

        const j: any =
            data.transitive_attested_function_call.transitive_attestation
        const taBytes: Uint8Array = ethers.decodeBase64(j)
        const ta: string = Buffer.from(taBytes).toString('hex');

        return {
            publicKey: `0x${publicKeyHex}`,
            transitiveAttestation: `0x${ta}`
        };
    } catch (e) {
        throw new Error(`Error loading EVM link data: ` + e);
    }
}

interface UserContract extends ethers.Contract {
    // @ts-ignore
    setEnclaveAttestedAppPubKey(publicKey:string): Promise<ethers.ContractTransactionResponse>;
    processTransitiveAttestedFunctionCall(ta: string): Promise<ethers.ContractTransactionResponse>;
}

describe("Local Test", function (): void {
    async function deployUser(): Promise<{ userContract: User }> {
        const contract: User = await hre.ethers.deployContract("User");
        return {userContract: contract};
    }

    it("Verify TA", async (): Promise<void> => {
        // given
        const taFile = process.env.TA_FILE;
        const evmLinkData: EVMLinkData = loadEVMLinkData(taFile);
        const publicKey: string = evmLinkData.publicKey;

        const {userContract} = await loadFixture(deployUser) as UserContract;

        await userContract.setEnclaveAttestedAppPubKey(publicKey);

        // when
        const ta: string = evmLinkData.transitiveAttestation;
        const tx: ethers.ContractTransactionResponse =
            await userContract.processTransitiveAttestedFunctionCall(ta);

        // then
        const expEvent = "AttestedFunctionCallOutput"
        const expEventArg = "Hello, World!"
        await expect(tx).to.emit(userContract, expEvent).withArgs(expEventArg);
        console.log("\tprocessTransitiveAttestedFunctionCall emitted %s(%s)", expEvent, expEventArg);
    })
});
