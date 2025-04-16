default: run

tmp:
	@mkdir -p tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

SRCS := $(shell find . -name '*.go')
DEPS := go.mod go.sum
TINYGO_IMAGE=tinygo/tinygo:0.34.0
CONTAINER_NAME=tinygo-builder
TINYGO_FLAGS=-no-debug -opt=z -scheduler=none

tmp/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM module..."
	@docker rm -f $(CONTAINER_NAME) > /dev/null 2>&1 || true
	@docker run  -dit --name $(CONTAINER_NAME) \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
	    $(TINYGO_IMAGE) sleep infinity > /dev/null
	@docker exec -u root $(CONTAINER_NAME) \
 		chown -R 1000:1000 /home/tinygo/.cache
	@docker exec -u root $(CONTAINER_NAME) \
 		chown -R 1000:1000 /home/tinygo/go
	@docker exec -w /src -u tinygo $(CONTAINER_NAME) \
 		tinygo build $(TINYGO_FLAGS) -o tmp/x.wasm -target=wasi ./...
	@docker stop $(CONTAINER_NAME) > /dev/null 2>&1 && docker rm $(CONTAINER_NAME) > /dev/null 2>&1 &

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

