test-hello-world-attest-fn-call:
	@echo "Running test-hello-world-attest-fn-call..."
	@make -C hello_world_attest_fn_call run > /dev/null 2>&1

test-hello-world-on-chain:
	@echo "Running test-hello-world-on-chain..."
	@make -C hello_world_on_chain test-local > /dev/null

test-error-handling:
	@echo "Running test-error-handling..."
	@make -C error_handling run-success > /dev/null 2>&1
	@make -C error_handling run-error > /dev/null 2>&1

test-coin-prices-from-coingecko:
	@echo "Running test-coin-prices-from-coingecko..."
	@make -C coin_prices_from_coingecko run > /dev/null 2>&1

test-twap-fn-call:
	@echo "Running test-twap-fn-call..."
	@make -C time_weighted_average_price/attest_fn_call init > /dev/null 2>&1
	@make -C time_weighted_average_price/attest_fn_call iteration > /dev/null 2>&1
	@make -C time_weighted_average_price/attest_fn_call twap > /dev/null 2>&1

test-twap-on-chain:
	@echo "Running test-twap-on-chain..."
	@make -C time_weighted_average_price/on_chain test-local > /dev/null

test-all: test-hello-world-attest-fn-call test-hello-world-on-chain test-error-handling test-coin-prices-from-coingecko test-twap-fn-call test-twap-on-chain
	@echo "All tests passed!"
