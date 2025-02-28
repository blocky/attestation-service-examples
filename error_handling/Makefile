tmp:
	@mkdir -p tmp

check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }
	
tmp/x.wasm: main.go $(wildcard as/*.go) | tmp check
	@echo "Building WASM module..."
	@docker run --rm \
        -v .:/src \
        -w /src \
        tinygo/tinygo:0.31.2 \
        tinygo build -o tmp/x.wasm -target=wasi main.go

build: tmp/x.wasm

run-success: build | check
	@echo "Running WASM module..."
	@cat successFunc.json | bky-as attest-fn-call >tmp/successout.json
	@echo "Output:"
	@jq -r '.function_calls[0].claims.output | @base64d | fromjson' tmp/successout.json

run-error: build | check
	@echo "Running WASM module..."
	@cat errorFunc.json | bky-as attest-fn-call >tmp/errorout.json
	@echo "Output:"
	@jq -r '.function_calls[0].claims.output | @base64d | fromjson' tmp/errorout.json

clean:
	@rm -rf tmp

