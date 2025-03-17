tmp:
	@mkdir tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

tmp/x.wasm: main.go $(wildcard price/*.go) | tmp check
	@echo "Building WASM module..."
	@docker run --rm \
        -v .:/src \
        -w /src \
        tinygo/tinygo:0.31.2 \
        tinygo build -o tmp/x.wasm -target=wasi ./...

.PHONY: build
build: tmp/x.wasm

.PHONY: prev
prev: check
	$(eval prev_ea := $(shell jq '.enclave_attested_application_public_key.enclave_attestation' tmp/prev.json | sed 's/\//\\\//g' ))
	$(eval prev_ta := $(shell jq '.transitive_attested_function_call.transitive_attestation' tmp/prev.json ))

.PHONY: init
init: check build
	@sed -e 's/VAR_EATTEST/null/' -e 's/VAR_TATTEST//' ./iteration-call.json.template > ./tmp/iteration-call.json
	@cat tmp/iteration-call.json | bky-as attest-fn-call  | jq . > tmp/prev.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d' tmp/prev.json | jq .

.PHONY: iteration
iteration: check prev build
	@sed \
		-e 's/VAR_TATTEST/"tAttest": ${prev_ta},/' \
		-e 's/VAR_EATTEST/${prev_ea}/' \
		iteration-call.json.template > tmp/iteration-call.json
	@cat tmp/iteration-call.json | bky-as attest-fn-call | jq . > tmp/prev.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/prev.json

.PHONY: twap
twap: check prev build
	@sed \
		-e 's/VAR_TATTEST/"tAttest": ${prev_ta},/' \
		-e 's/VAR_EATTEST/${prev_ea}/' \
		twap-call.json.template > tmp/twap-call.json
	@cat tmp/twap-call.json | bky-as attest-fn-call | jq . > tmp/twap.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/twap.json

.PHONY: clean
clean:
	@rm -rf tmp

