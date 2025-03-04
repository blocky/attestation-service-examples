.PHONY: test-local
test-local:
	npx hardhat test --grep "Local"

.PHONY: copy-from-twap
copy-from-twap:
	@mkdir -p inputs/
	@cp ../twap/tmp/twap.json inputs/

.PHONY: very-clean
very-clean:
	@rm -rf artifacts cache node_modules typechain-types
