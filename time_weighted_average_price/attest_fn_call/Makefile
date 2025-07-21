tmp:
	@mkdir tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v bky-c >/dev/null 2>&1 || { echo >&2 "bky-c not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

SRCS := $(shell find . -name '*.go')
DEPS := go.mod go.sum

tmp/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM module..."
	@mkdir -p tmp
	@bky-c build . tmp/x.wasm

.PHONY: build
build: tmp/x.wasm

.PHONY: prev
prev: check
	$(eval prev_ea := $(shell jq -r '.enclave_attested_application_public_key.enclave_attestation' tmp/prev.json))
	$(eval prev_ta := $(shell jq -r '.transitive_attested_function_call.transitive_attestation' tmp/prev.json ))

.PHONY: init
init: check build
	@sed \
		-e 's|PREV_ENCLAVE_ATTESTATION|null|' \
		-e 's|PREV_TRANSITIVE_ATTESTATION|null|' \
		./iteration-call.json.template > ./tmp/iteration-call.json
	@cat tmp/iteration-call.json | bky-as attest-fn-call  | jq . > tmp/prev.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d' tmp/prev.json | jq .

.PHONY: iteration
iteration: check prev build
	@sed \
		-e 's|PREV_ENCLAVE_ATTESTATION|"${prev_ea}"|' \
		-e 's|PREV_TRANSITIVE_ATTESTATION"|${prev_ta}"|' \
		iteration-call.json.template > tmp/iteration-call.json
	@cat tmp/iteration-call.json | bky-as attest-fn-call | jq . > tmp/prev.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/prev.json

.PHONY: twap
twap: check prev build
	@sed \
		-e 's|PREV_ENCLAVE_ATTESTATION|"${prev_ea}"|' \
		-e 's|PREV_TRANSITIVE_ATTESTATION|"${prev_ta}"|' \
		twap-call.json.template > tmp/twap-call.json
	@cat tmp/twap-call.json | bky-as attest-fn-call | jq . > tmp/twap.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/twap.json

.PHONY: clean
clean:
	@rm -rf tmp
