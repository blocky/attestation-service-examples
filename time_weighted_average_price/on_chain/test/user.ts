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
    processTransitivelyAttestedResult(publicKey: string, ta: string): Promise<ethers.ContractTransactionResponse>;
}

describe("Local Tests", function () {
    async function deployUser(): Promise<{ userContract: User }> {
        const contract: User = await hre.ethers.deployContract("User");
        return {userContract: contract};
    }

    it("Verify attested TWAP in User contract", async () => {
        // given
        const taFile = process.env.TA_FILE || "../inputs/twap.json";
        const evmLinkData: EVMLinkData = loadEVMLinkData(taFile);
        const {userContract} = await loadFixture(deployUser) as UserContract;

        function anyFloat(x: unknown): boolean {
            const n = Number(x);
            return !Number.isNaN(n) && Number.isFinite(n);
        }

        // when
        const tx: ethers.ContractTransactionResponse =
            await userContract.processTransitivelyAttestedResult(
                evmLinkData.publicKey,
                evmLinkData.transitiveAttestation,
            )

        // then
        // assert the transaction emitted an event with a float number
        // (time weighted average price) as an argument
        await expect(tx).to.emit(
            userContract,
            'TWAP'
        ).withArgs(anyFloat)
    })
});
