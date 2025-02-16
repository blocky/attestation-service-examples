import hre from "hardhat";
import {loadFixture} from "@nomicfoundation/hardhat-toolbox/network-helpers";
import {expect} from "chai";
import {ethers} from "ethers";

const fs = require("fs")
const path = require("path")

const loadEVMLinkData = () => {
    try {
        const dir = path.resolve( __dirname, "../inputs/twap.json");
        const file = fs.readFileSync(dir);

        const data = JSON.parse(file);

        const k = data.enclave_attested_application_public_key.public_key.data
        const pubKeyBytes = ethers.decodeBase64(k)
        const publicKeyHex = Buffer.from(pubKeyBytes).toString('hex');

        const taBytes = ethers.decodeBase64(data.function_calls[0].transitive_attestation)
        const ta = Buffer.from(taBytes).toString('utf-8');



        return {
            publicKey: `0x${publicKeyHex}`,
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
        // given
        const evmLinkData = loadEVMLinkData();
        const publicKey = evmLinkData.publicKey;

        const {userContract} = await loadFixture(deployUser);
        await userContract.setTASigningKeyAddress(publicKey as any);

        // when
        const ta = evmLinkData.transitiveAttestation;
        const tx = await userContract.processAttestedFnCallClaims(ta as any)

        // then
        await expect(tx).to.emit(
            userContract,
            'AttestedFunctionCallOutput'
        )
    })
});



