default: run

tmp:
	@mkdir -p tmp
	
tmp/x.wasm: main.go | tmp
	tinygo build -o tmp/x.wasm -target=wasi main.go

build: tmp/x.wasm

run: build
	cat fn-call.json | ./bky-as attest-fn-call  | jq -r .function_calls[0].claims.output | base64 -d | jq .

clean:
	rm -rf tmp

