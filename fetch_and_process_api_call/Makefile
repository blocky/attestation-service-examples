default: run

tmp:
	@mkdir -p tmp
	
tmp/x.wasm: main.go | tmp
	tinygo build -o tmp/x.wasm -target=wasi main.go
	#docker run --rm -v .:/src -w /src tinygo/tinygo:0.35.0 tinygo build -o tmp/x.wasm -target=wasi main.go

build: tmp/x.wasm

run: build
	export PATH=$$PATH:$$PWD && \
	cat fn-call.json | bky-as attest-fn-call  | jq -r '.function_calls[0].claims.output | @base64d | fromjson'

clean:
	rm -rf tmp

