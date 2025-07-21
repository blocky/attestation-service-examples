default: run-success

tmp:
	@mkdir -p tmp

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

.PHONY: run-success
run-success: check build
	@echo "Running WASM module..."
	@cat successFunc.json | bky-as attest-fn-call >tmp/out-success.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out-success.json

.PHONY: run-error
run-error: check build
	@echo "Running WASM module..."
	@cat errorFunc.json | bky-as attest-fn-call >tmp/out-error.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out-error.json

.PHONY: run-panic
run-panic: check build
	@echo "Running WASM module..."
	@cat panicFunc.json | bky-as attest-fn-call >tmp/out-panic.json

.PHONY: clean
clean:
	@rm -rf tmp

