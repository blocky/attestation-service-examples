default: match-winner

tmp:
	@mkdir -p tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v bky-c >/dev/null 2>&1 || { echo >&2 "bky-c not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

SRCS := $(shell find . -name '*.go')
DEPS := go.mod go.sum

tmp/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM module..."
	@bky-c build . tmp/x.wasm

.PHONY: build
build: tmp/x.wasm

.PHONY: match-winner
match-winner: check build
	@echo "Running WASM module..."
	@sed -e 's/{{ .YOUR_RIMBLE_MATCH_DATE }}/$(RIMBLE_MATCH_DATE)/' \
	     -e 's/{{ .YOUR_RIMBLE_MATCH_ID }}/$(RIMBLE_MATCH_ID)/' \
	     -e 's/{{ .YOUR_RIMBLE_API_KEY }}/$(RIMBLE_API_KEY)/' \
	     ./match-winner.json.template > ./tmp/match-winner.json
	@cat ./tmp/match-winner.json | bky-as attest-fn-call >tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out.json

.PHONY: team-kill-diff
team-kill-diff: check build
	@echo "Running WASM module..."
	@sed -e 's/{{ .YOUR_RIMBLE_MATCH_DATE }}/$(RIMBLE_MATCH_DATE)/' \
	     -e 's/{{ .YOUR_RIMBLE_MATCH_ID }}/$(RIMBLE_MATCH_ID)/' \
	     -e 's/{{ .YOUR_RIMBLE_API_KEY }}/$(RIMBLE_API_KEY)/' \
	     ./team-kill-diff.json.template > ./tmp/team-kill-diff.json
	@cat ./tmp/team-kill-diff.json | bky-as attest-fn-call >tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/out.json

.PHONY: fetch-match-id-and-date
fetch-match-id-and-date:
	@echo "Fetching recent match ID and date..."
	@curl -s --location 'https://rimbleanalytics.com/raw/csgo/completed-matches/' \
		--header 'x-api-key: $(RIMBLE_API_KEY)' \
		| jq -r '.[0] | "Match ID: \(.matchid)\nDate: \(.date)"'

.PHONY: clean
clean:
	@rm -rf tmp

.PHONY: test-rimble
test-rimble:
	go test -v ./rimble/...
