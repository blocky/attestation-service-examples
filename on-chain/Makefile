.PHONY: test-local
test-local:
	npx hardhat test --grep "Local"

.PHONY: compile
compile:
	npx hardhat compile

.PHONY: deploy
deploy: compile
	npx hardhat run scripts/deployUser.ts --network base-sepolia
	npx hardhat verify --network base-sepolia $$(cat .user_deployed_address)

