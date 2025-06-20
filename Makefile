.PHONY: test
test:
	@go test -C test . -count=1

.PHONY: test-live
test-live:
	@go test -C test ./live -count=1

.PHONY: pre-pr
pre-pr: test
