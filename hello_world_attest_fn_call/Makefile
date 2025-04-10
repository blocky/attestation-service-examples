default: run

tmp:
	@mkdir -p tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

init:
	@install -d -o 1000 -g 1000 ${HOME}/.cache/bky-as/tinygo-dot-cache
	@install -d -o 1000 -g 1000 ${HOME}/.cache/bky-as/tinygo-go-path

build: 
	@echo "Building WASM **module..."
	@docker run --rm \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
		tinygo/tinygo:0.31.2 \
		tinygo build -x -o tmp/x.wasm -target=wasi ./...

.PHONY: run
run: check build
	@echo "Running WASM module..."
	@cat fn-call.json | bky-as attest-fn-call >tmp/out.json
	@echo "Log:"
	@jq -r '.transitive_attested_function_call.logs | @base64d' tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d ' tmp/out.json

.PHONY: clean
clean:
	@rm -rf tmp

