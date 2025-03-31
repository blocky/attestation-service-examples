import hre from "hardhat";
import {loadFixture} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import {expect} from "chai";
import {ethers} from "ethers";

const fs = require("fs")
const path = require("path")

function loadEVMLinkData(jsonPath: string) {
    try {
        const dir = path.resolve( __dirname, jsonPath);
        const file = fs.readFileSync(dir);

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
        console.log(`e`, e)
    }
}

const loadUserDeployedAddress = () => {
    try {
        const dir = path.resolve(
            __dirname,
            "../deployments/user_deployed_address"
        )
        const file = fs.readFileSync(dir, "utf8")

        return file.toString()
    } catch (e) {
        console.log(`e`, e)
    }
}

const loadUserContractABI = () => {
    try {
        const dir = path.resolve(
            __dirname,
            "../artifacts/contracts/User.sol/User.json"
        )
        const file = fs.readFileSync(dir, "utf8")
        const json = JSON.parse(file)
        return json.abi
    } catch (e) {
        console.log(`e`, e)
    }
}

describe("Local Test", function () {
    async function deployUser() {
        const contract = await hre.ethers.deployContract("User");
        return {userContract: contract};
    }

    it("Verify TA and parse Result w/success", async () => {
        // given
        const evmLinkData = loadEVMLinkData("../inputs/out-success.json");
        const publicKey = evmLinkData.publicKey;

        const {userContract} = await loadFixture(deployUser);
        await userContract.setTASigningKeyAddress(publicKey as any);

        const ta = evmLinkData.transitiveAttestation;

        // when
        const tx = await userContract.verifyAttestedFnCallClaims(ta as any)

        // then
        await expect(tx).to.emit(
            userContract,
            'ResultValue'
        ).withArgs("{\"number\":42}")
    })

    it("Verify TA and parse Result w/error", async () => {
        // given
        const evmLinkData = loadEVMLinkData("../inputs/out-error.json");
        const publicKey = evmLinkData.publicKey;

        const {userContract} = await loadFixture(deployUser);
        await userContract.setTASigningKeyAddress(publicKey as any);

        const ta = evmLinkData.transitiveAttestation;

        // when/then
        await expect(
            userContract.verifyAttestedFnCallClaims(ta as any)
        ).to.be.revertedWith("expected error")
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

    const evmLinkData = loadEVMLinkData("../inputs/out-success.json");

    it("Set signing key", async () => {
        const publicKey = evmLinkData.publicKey;
        const tx = await userContract.setTASigningKeyAddress(publicKey as any);
        await tx.wait()
    })

    it("Verify TA", async () => {
        const ta = evmLinkData.transitiveAttestation;
        const tx = await userContract.verifyAttestedFnCallClaims(ta)
        // poll instead of tx.wait() to get the lowest possible delay
        for (; ;) {
            const txReceipt = await provider.getTransactionReceipt(tx.hash);
            if (txReceipt && txReceipt.blockNumber) {
                break
            }
        }
    });
});


