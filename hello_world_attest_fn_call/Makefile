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

FUNCTION := "hello-world"

run: build | check
	@echo "Running WASM module..."
	@cat $(FUNCTION)-call.json \
	| bky-as attest-fn-call > tmp/out.json
	@echo "Result:"
	@jq -r '.function_calls[0].claims.output | @base64d | fromjson' tmp/out.json
	@LOG=$$(jq -r '.function_calls[0].logs' tmp/out.json); \
	if [ "$$LOG" != "null" ]; \
	then echo "Logs:"; jq -r '.function_calls[0].logs | @base64d' tmp/out.json; \
	fi

hello-world:
	@$(MAKE) --no-print-directory run FUNCTION="hello-world"

hello-error:
	@$(MAKE) --no-print-directory run FUNCTION="hello-error"

clean:
	@rm -rf tmp

