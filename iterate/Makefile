prev_ta := $(shell jq '.function_calls[0].transitive_attestation' tmp/prev.json) 

default: run

tmp: 
	mkdir tmp
	
build: tmp
	tinygo build -o tmp/x.wasm -target=wasi main.go

init: 
	@sed -e "s/TATTEST//" fn-call.json.template > tmp/fn-call.json
	@cat tmp/fn-call.json | ./bky-as attest-fn-call  | jq . > tmp/prev.json 
	@jq -r '.function_calls[0].claims.output | @base64d' tmp/prev.json | jq .

run: 
	@sed -e "s/TATTEST/\"tAttest\": \"${prev_ta}\",/" -e 's/ \",/",/' fn-call.json.template > tmp/fn-call.json
	@cat tmp/fn-call.json | ./bky-as attest-fn-call  | jq . > tmp/prev.json 
	@jq -r '.function_calls[0].claims.output | @base64d' tmp/prev.json | jq .

clean:
	rm -rf tmp

