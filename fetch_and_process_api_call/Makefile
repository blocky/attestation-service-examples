default: run

tmp:
	@mkdir -p tmp
	
tmp/x.wasm: main.go $(wildcard as/*.go) | tmp
	@echo "Building WASM module..."
	@tinygo build -o tmp/x.wasm -target=wasi main.go

build: tmp/x.wasm

run: build
	@echo "Running WASM module..."
	@export PATH=$$PATH:$(pwd) && \
	cat fn-call.json \
	| bky-as attest-fn-call \
	| jq -r '.function_calls[0].claims.output | @base64d | fromjson'

clean:
	rm -rf tmp

