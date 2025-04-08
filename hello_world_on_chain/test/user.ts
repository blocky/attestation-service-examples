// noinspection DuplicatedCode

import hre from "hardhat";
import {loadFixture} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import {expect} from "chai";
import {ethers} from "ethers";
import fs from "fs";
import path from "path";

function loadEVMLinkData(
    jsonPath: string,
): { publicKey: string, transitiveAttestation: string} {
    try {
        const dir = path.resolve( __dirname, jsonPath);
        const file = fs.readFileSync(dir, "utf8");

        const data = JSON.parse(file);

        const k = data.enclave_attested_application_public_key.claims.public_key.data
        const pubKeyBytes = ethers.decodeBase64(k)
        const publicKeyHex = Buffer.from(pubKeyBytes).toString('hex');

        const taBytes = ethers.decodeBase64(data.transitive_attested_function_call.transitive_attestation)
        const ta = Buffer.from(taBytes).toString('utf-8');

        return {
            publicKey: `0x${publicKeyHex}`,
            transitiveAttestation: ta
        };
    } catch (e) {
        throw new Error(`Error loading EVM link data: ` + e);
    }
}

const loadUserDeployedAddress: () => string  = (
) : string  => {
    try {
        const dir = path.resolve(
            __dirname,
            "../deployments/user_deployed_address"
        )
        const file = fs.readFileSync(dir, "utf8")

        return file.toString()
    } catch (e) {
        throw new Error(`loading user deployed address: ` + e);
    }
}

const loadUserContractABI: () => any = (
) : any =>  {
    try {
        const dir = path.resolve(
            __dirname,
            "../artifacts/contracts/User.sol/User.json"
        )
        const file = fs.readFileSync(dir, "utf8")
        const json = JSON.parse(file)
        return json.abi
    } catch (e) {
        throw new Error(`loading user contract ABI: ` + e);
    }
}

interface UserContract extends ethers.Contract {
    processTAHelloWorld(publicKey: any, ta: any): Promise<ethers.ContractTransactionResponse>;
}

describe("Local Test", function () {
    async function deployUser() {
        const contract = await hre.ethers.deployContract("User");
        return {userContract: contract};
    }

    it("Verify TA", async () => {
        // given
        const evmLinkData = loadEVMLinkData("../inputs/out.json");
        const publicKey = evmLinkData.publicKey;

        const {userContract} = await loadFixture(deployUser);

        // when
        const ta = evmLinkData.transitiveAttestation;
        const tx = await userContract.processTAHelloWorld(
            publicKey as any,
            ta as any
        )

        // then
        await expect(tx).to.emit(
            userContract,
            'AttestedFunctionCallOutput'
        ).withArgs("Hello, World!")
    })
});

describe("Base Sepolia Tests", function () {
    const url = 'https://sepolia.base.org';
    const provider = new ethers.JsonRpcProvider(url);
    const privateKey = process.env.WALLET_KEY as string
    const signer = new ethers.Wallet(privateKey, provider)

    const userContract = new ethers.Contract(
        loadUserDeployedAddress(),
        loadUserContractABI(),
        signer
    );

    const evmLinkData = loadEVMLinkData("../inputs/out.json");

    it("Verify TA", async () => {
        const publicKey = evmLinkData.publicKey;
        const ta = evmLinkData.transitiveAttestation;
        const tx = await userContract.processTAHelloWorld(
            publicKey as any,
            ta as any
        )
        // poll instead of tx.wait() to get the lowest possible delay
        for (; ;) {
            const txReceipt = await provider.getTransactionReceipt(tx.hash);
            if (txReceipt && txReceipt.blockNumber) {
                break
            }
        }
    });
});


