.PHONY: test-local
test-local:
	npx hardhat test --grep "Local"

.PHONY: copy-from-iterate
copy-from-iterate:
	mkdir -p inputs/
	cp ../iterate/tmp/prev.json inputs/

.PHONY: very-clean
very-clean:
	@rm -rf artifacts cache node_modules typechain-types
