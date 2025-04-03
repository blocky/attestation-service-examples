run-hello-world-attest-fn-call:
	@echo "Running run-hello-world-attest-fn-call..."
	@make -C hello_world_attest_fn_call run > /dev/null 2>&1

run-hello-world-on-chain:
	@echo "Running run-hello-world-on-chain..."
	@make -C hello_world_on_chain test-local > /dev/null

run-error-handling-attest-fn-call:
	@echo "Running run-error-handling-attest-fn-call..."
	@make -C error_handling_attest_fn_call run-success > /dev/null 2>&1
	@make -C error_handling_attest_fn_call run-error > /dev/null 2>&1
	@make -C error_handling_attest_fn_call run-panic > /dev/null 2>&1

run-error-handling-on-chain:
	@echo "Running run-error-handling-on-chain..."
	@make -C error_handling_on_chain test-local > /dev/null

run-coin-prices-from-coingecko:
	@echo "Running run-coin-prices-from-coingecko..."
	@make -C coin_prices_from_coingecko run > /dev/null 2>&1

run-esports-data-from-pandascore:
	@echo "Running run-esports-data-from-pandascore..."
	@make -C esports_data_from_pandascore run > /dev/null 2>&1

run-shipment_tracking_with_dhl:
	@echo "Running run-shipment_tracking_with_dhl..."
	@make -C shipment_tracking_with_dhl run > /dev/null 2>&1

run-twap-fn-call:
	@echo "Running run-twap-fn-call..."
	@make -C time_weighted_average_price/attest_fn_call init > /dev/null 2>&1
	@make -C time_weighted_average_price/attest_fn_call iteration > /dev/null 2>&1
	@make -C time_weighted_average_price/attest_fn_call twap > /dev/null 2>&1

run-twap-on-chain:
	@echo "Running run-twap-on-chain..."
	@make -C time_weighted_average_price/on_chain test-local > /dev/null

run-random:
	@echo "Running random..."
	@make -C random clean run 2> /dev/null

run-params-and-secrets:
	@echo "Running params-and-secrets..."
	@make -C params_and_secrets run > /dev/null 2>&1
	@make -C params_and_secrets run-error > /dev/null 2>&1

run-all: \
	run-hello-world-attest-fn-call \
	run-hello-world-on-chain \
	run-error-handling-attest-fn-call \
	run-error-handling-on-chain \
	run-coin-prices-from-coingecko \
	run-esports-data-from-pandascore \
	run-shipment_tracking_with_dhl \
	run-twap-fn-call  \
	run-twap-on-chain \
	run-random \
	run-params-and-secrets
	@echo "All tests passed!"
