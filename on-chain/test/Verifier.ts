import hre from "hardhat";
import {loadFixture} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import {expect} from "chai";
import {ethers} from "ethers";

require('dotenv').config();

const fs = require("fs")
const path = require("path")

const loadEVMLinkData = () => {
    try {
        const dir = path.resolve( __dirname, "../inputs/prev.json");
        const file = fs.readFileSync(dir);

        const data = JSON.parse(file);

        const bytes = ethers.decodeBase64(data.function_calls[0].transitive_attestation)
        const ta = Buffer.from(bytes).toString('utf-8');
        return {
            publicKey: `0x${data.enclave_attested_application_public_key.public_key.data}`,
            transitiveAttestation: ta
        };
    } catch (e) {
        console.log(`e`, e)
    }
}

describe("Local Tests", function () {
    async function deployUser() {
        const contract = await hre.ethers.deployContract("User");
        return {userContract: contract};
    }

    it("process TA", async () => {
        const {userContract} = await loadFixture(deployUser);
        const evmLinkData = loadEVMLinkData();

        const publicKey = evmLinkData.publicKey;
        await userContract.setTASigningKeyAddress(publicKey as any);

        const ta = evmLinkData.transitiveAttestation;
        const tx = await userContract.processAttestedAPICallClaims(ta)

        // then
        await expect(tx).to.emit(
            userContract,
            'AttestedFunctionCallOutput'
        )
    })
});



