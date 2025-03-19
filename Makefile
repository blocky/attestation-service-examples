run-hello-world-attest-fn-call:
	@echo "Running run-hello-world-attest-fn-call..."
	@make -C hello_world_attest_fn_call clean run > /dev/null 2>&1

run-hello-world-on-chain:
	@echo "Running run-hello-world-on-chain..."
	@make -C hello_world_on_chain test-local > /dev/null

run-error-handling:
	@echo "Running run-error-handling..."
	@make -C error_handling clean run-success > /dev/null 2>&1
	@make -C error_handling clean run-error > /dev/null 2>&1
	@make -C error_handling clean run-panic > /dev/null 2>&1

run-coin-prices-from-coingecko:
	@echo "Running run-coin-prices-from-coingecko..."
	@make -C coin_prices_from_coingecko clean run > /dev/null 2>&1

run-twap-fn-call:
	@echo "Running run-twap-fn-call..."
	@make -C time_weighted_average_price/attest_fn_call clean init > /dev/null 2>&1
	@make -C time_weighted_average_price/attest_fn_call iteration > /dev/null 2>&1
	@make -C time_weighted_average_price/attest_fn_call twap > /dev/null 2>&1

run-esports-data-from-pandascore:
	@echo "Running run-esports-data-from-pandascore..."
	@make -C esports_data_from_pandascore clean run > /dev/null 2>&1

run-twap-on-chain:
	@echo "Running run-twap-on-chain..."
	@make -C time_weighted_average_price/on_chain test-local > /dev/null

run-all: run-hello-world-attest-fn-call run-hello-world-on-chain run-error-handling run-coin-prices-from-coingecko run-twap-fn-call run-esports-data-from-pandascore run-twap-on-chain
	@echo "All tests passed!"
