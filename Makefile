run-hello-world-attest-fn-call:
	@echo "Running run-hello-world-attest-fn-call..."
	@cd hello_world_attest_fn_call && make run 2> /dev/null

run-hello-world-on-chain:
	@echo "Running run-hello-world-on-chain..."
	@cd hello_world_on_chain && make test-local 2> /dev/null

run-error-handling-attest-fn-call:
	@echo "Running run-error-handling-attest-fn-call..."
	@cd error_handling_attest_fn_call && make run-success 2> /dev/null
	@cd error_handling_attest_fn_call && make run-error 2> /dev/null
	@cd error_handling_attest_fn_call && make run-panic 2> /dev/null

run-error-handling-on-chain:
	@echo "Running run-error-handling-on-chain..."
	@cd error_handling_on_chain && @make  test-local 2> /dev/null

run-coin-prices-from-coingecko:
	@echo "Running run-coin-prices-from-coingecko..."
	@cd coin_prices_from_coingecko && make run 2> /dev/null

run-esports-data-from-pandascore:
	@echo "Running run-esports-data-from-pandascore..."
	@cd esports_data_from_pandascore && make run 2> /dev/null

run-shipment_tracking_with_dhl:
	@echo "Running run-shipment_tracking_with_dhl..."
	@cd shipment_tracking_with_dhl && make run 2> /dev/null

run-twap-fn-call:
	@echo "Running run-twap-fn-call..."
	@cd time_weighted_average_price/attest_fn_call && make init 2> /dev/null
	@cd time_weighted_average_price/attest_fn_call && make iteration 2> /dev/null
	@cd time_weighted_average_price/attest_fn_call && make twap 2> /dev/null

run-twap-on-chain:
	@echo "Running run-twap-on-chain..."
	@cd time_weighted_average_price/on_chain && make test-local 2> /dev/null

run-random:
	@echo "Running random..."
	@cd random &&  make run 2> /dev/null

run-time:
	@echo "Running time..."
	@cd time &&  make run 2> /dev/null

run-params-and-secrets:
	@echo "Running params-and-secrets..."
	@cd params_and_secrets && make run 2> /dev/null
	@cd params_and_secrets && make run-error 2> /dev/null

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
	run-time \
	run-params-and-secrets
	@echo "All tests passed!"
