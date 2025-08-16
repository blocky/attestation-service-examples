# test-local needs TA_FILE env var set to path of out.json
.PHONY: test-local
test-local:
	npx hardhat test --grep "User contract"

.PHONY: very-clean
very-clean:
	@rm -rf artifacts cache node_modules typechain-types tmp
