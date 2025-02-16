prev_ta := $(shell jq '.function_calls[0].transitive_attestation' tmp/prev.json )
prev_ea := $(shell jq '.enclave_attested_application_public_key.enclave_attestation' tmp/prev.json | sed 's/\//\\\//g' )

default: run

tmp: 
	@mkdir tmp

check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

tmp/x.wasm: main.go $(wildcard as/*.go) $(wildcard price/*.go) | tmp check
	@echo "Building WASM module..."
	@docker run --rm \
        -v .:/src \
        -w /src \
        tinygo/tinygo:0.31.2 \
        tinygo build -o tmp/x.wasm -target=wasi main.go

build: tmp/x.wasm

init: build
	@sed -e 's/VAR_EATTEST/null/' -e 's/VAR_TATTEST//' ./iteration-call.json.template > ./tmp/iteration-call.json
	@cat tmp/iteration-call.json | ./bky-as attest-fn-call  | jq . > tmp/prev.json
	@jq -r '.function_calls[0].claims.output | @base64d' tmp/prev.json | jq .

echo:
	@echo "prev_ea: ${prev_ea}"

iteration: build | check
	@sed \
		-e 's/VAR_TATTEST/"tAttest": ${prev_ta},/' \
		-e 's/VAR_EATTEST/${prev_ea}/' \
		iteration-call.json.template > tmp/iteration-call.json
	@cat tmp/iteration-call.json | ./bky-as attest-fn-call | jq . > tmp/prev.json
	@jq -r '.function_calls[0].claims.output | @base64d' tmp/prev.json | jq .
	@jq -r '.function_calls[0].logs | @base64d' tmp/prev.json

twap: build | check
	@sed \
		-e 's/VAR_TATTEST/"tAttest": ${prev_ta},/' \
		-e 's/VAR_EATTEST/${prev_ea}/' \
		twap-call.json.template > tmp/twap-call.json
	@cat tmp/twap-call.json | ./bky-as attest-fn-call | jq . > tmp/twap.json
	@jq -r '.function_calls[0].claims.output | @base64d' tmp/twap.json | jq .
	@jq -r '.function_calls[0].logs | @base64d' tmp/twap.json

clean:
	@rm -rf tmp

