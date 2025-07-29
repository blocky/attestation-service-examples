default: run

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v bky-c >/dev/null 2>&1 || { echo >&2 "bky-c not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

SRCS := $(shell find . -name '*.go')
DEPS := go.mod go.sum

main.wasm: $(SRCS) $(DEPS) | check
	@echo "Building WASM module..."
	@mkdir -p -m 777 tmp
	@bky-c build . tmp/x.wasm

out.json: main.wasm fn-call.json | check
	@cat fn-call.json | bky-as attest-fn-call >out.json

.PHONY: run
run: out.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d ' out.json

.PHONY: clean
clean:
	@rm -f main.wasm out.json
