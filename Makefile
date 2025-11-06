# Optionally set additional flags for `go test` command. If specifying more than
# one flag, separate them with spaces.
GO_TEST_FLAGS :=

.PHONY: test
test:
	@go test -C test . -count=1 $(GO_TEST_FLAGS)

.PHONY: test-live
test-live:
	@go test -C test ./live -count=1 $(GO_TEST_FLAGS)

.PHONY: pre-pr
pre-pr: test
