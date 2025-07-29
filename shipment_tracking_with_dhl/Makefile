default: run

tmp:
	@mkdir -p -m 777 tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v bky-c >/dev/null 2>&1 || { echo >&2 "bky-c not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

SRCS := $(shell find . -name '*.go')
DEPS := go.mod go.sum
TINYGO_IMAGE=tinygo/tinygo:0.34.0
CONTAINER_NAME=tinygo-builder
TINYGO_FLAGS=-no-debug -opt=z -scheduler=none

tmp/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM module..."
	@bky-c build . tmp/x.wasm

.PHONY: build
build: tmp/x.wasm

.PHONY: run
run: check build
	@echo "Running WASM module..."
	@cat fn-call.json | bky-as attest-fn-call >tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out.json

.PHONY: clean
clean:
	@rm -rf tmp

