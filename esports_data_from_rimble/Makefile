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

tmp/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM module..."
	@docker run --rm \
        -v .:/src \
        -w /src \
        tinygo/tinygo:0.31.2 \
        tinygo build -o tmp/x.wasm -target=wasi ./...

.PHONY: build
build: tmp/x.wasm

.PHONY: match-winner
match-winner: check build
	@echo "Running WASM module..."
	@sed -e 's/Rimble API key/$(RIMBLE_API_KEY)/' ./match-winner.json.template > ./tmp/match-winner.json
	@cat ./tmp/match-winner.json | bky-as attest-fn-call >tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out.json

.PHONY: team-kill-diff
team-kill-diff: check build
	@echo "Running WASM module..."
	@sed -e 's/Rimble API key/$(RIMBLE_API_KEY)/' ./team-kill-diff.json.template > ./tmp/team-kill-diff.json
	@cat ./tmp/team-kill-diff.json | bky-as attest-fn-call >tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out.json

.PHONY: clean
clean:
	@rm -rf tmp

.PHONY: update-rimble-test-data
update-rimble-test-data:
	@curl -s "https://rimbleanalytics.com/raw/csgo/match-status/?matchid=2379357&date=2025-02-18" \
		-H "x-api-key: $(RIMBLE_API_KEY)" \
		> rimble/testdata/match_data.json

.PHONY: test-rimble
test-rimble:
	go test -v ./rimble/...
