default: run

tmp:
	@mkdir -p tmp

.PHONY: check
check:
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker not found in PATH. Aborting."; exit 1; }
	@command -v bky-as >/dev/null 2>&1 || { echo >&2 "bky-as not found in PATH. Aborting."; exit 1; }
	@command -v jq >/dev/null 2>&1 || { echo >&2 "jq not found in PATH. Aborting."; exit 1; }

init:
	@install -d -o 1000 -g 1000 ${HOME}/.cache/bky-as/tinygo-dot-cache
	@install -d -o 1000 -g 1000 ${HOME}/.cache/bky-as/tinygo-go-path

SRCS := $(shell find . -name '*.go')
DEPS := go.mod go.sum
TINYGO_IMAGE=tinygo/tinygo:0.34.0
CONTAINER_NAME=tinygo-builder
TINYGO_FLAGS=-no-debug -opt=2 -scheduler=none

# Working version, ~minute cold start, ~6s with caching
tmp1/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM **module..."
	-docker rm -f tinygo-builder 2>/dev/null || true
	-docker run -dit --name tinygo-builder \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
	    tinygo/tinygo:0.31.2 sleep infinity
	-docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/.cache
	-docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/go
	-docker exec -w /src  -u tinygo tinygo-builder \
 		tinygo build -x -o tmp/x.wasm -target=wasi ./...
	-docker stop tinygo-builder > /dev/null 2>&1 && docker rm tinygo-builder > /dev/null 2>&1 &

# Working but after clean consistently 40-60s
tmp2/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM **module..."
	@docker run --rm -u root \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
		tinygo/tinygo:0.31.2 \
		sh -c "chown -R 1000:1000 /home/tinygo/.cache && chown -R 1000:1000 /home/tinygo/go && tinygo build -x -o tmp/x.wasm -target=wasi ./..."

# tinygo 0.33.0 ~5s, 0.34.0 ~6s, remove -x from build: ~4.5s
tmp3/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM **module..."
	-@docker rm -f tinygo-builder 2>/dev/null || true
	-@docker run -dit --name tinygo-builder \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
	    tinygo/tinygo:0.33.0 sleep infinity
	-@docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/.cache
	-@docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/go
	-@docker exec -w /src  -u tinygo tinygo-builder \
 		tinygo build -o tmp/x.wasm -target=wasi ./...
	-@docker stop tinygo-builder > /dev/null 2>&1 && docker rm tinygo-builder > /dev/null 2>&1 &

# cold 7s, cached 4.5s
tmp4/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM **module..."
	-docker rm -f tinygo-builder >/dev/null || true
	-docker run -dit --name tinygo-builder \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
	    tinygo/tinygo:0.33.0 sleep infinity > /dev/null
	-docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/.cache
	-docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/go
	-docker exec -w /src  -u tinygo tinygo-builder \
 		tinygo build -o tmp/x.wasm -target=wasi ./...
	-docker stop tinygo-builder > /dev/null 2>&1 && docker rm tinygo-builder > /dev/null 2>&1 &

# 6s cold, 4.1 cached
tmp5/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM **module..."
	-@docker rm -f tinygo-builder > /dev/null 2>&1 || true
	-@docker run --init --mount type=tmpfs,destination=${PWD}/tmp -di --name tinygo-builder \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
	    tinygo/tinygo:0.34.0 sleep infinity > /dev/null
	-@docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/.cache
	-@docker exec -u root tinygo-builder \
 		chown -R 1000:1000 /home/tinygo/go
	-@docker exec -w /src  -u tinygo tinygo-builder \
 		tinygo build --no-debug -opt=z -scheduler=none -o tmp/x.wasm -target=wasi ./...
	-@docker stop tinygo-builder > /dev/null 2>&1 && docker rm tinygo-builder > /dev/null 2>&1 &

tmp/x.wasm: $(SRCS) $(DEPS) | tmp check
	@echo "Building WASM **module..."
	-@docker rm -f $(CONTAINER_NAME) > /dev/null 2>&1 || true
	-@docker run -di --name $(CONTAINER_NAME) \
		-v .:/src \
		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
		-w /src \
	    $(TINYGO_IMAGE) sleep infinity > /dev/null
	-@docker exec -u root $(CONTAINER_NAME) \
 		chown -R 1000:1000 /home/tinygo/.cache
	-@docker exec -u root $(CONTAINER_NAME) \
 		chown -R 1000:1000 /home/tinygo/go
	-@docker exec -w /src -u tinygo $(CONTAINER_NAME) \
 		tinygo build $(TINYGO_FLAGS) -o tmp/x.wasm -target=wasi ./...
	-@docker stop $(CONTAINER_NAME) > /dev/null 2>&1 && docker rm $(CONTAINER_NAME) > /dev/null 2>&1 &

.PHONY: build
build: tmp/x.wasm

#	@docker run --rm -u root \
#		-v .:/src \
#		-v ${HOME}/.cache/bky-as/tinygo-dot-cache:/home/tinygo/.cache \
#		-v ${HOME}/.cache/bky-as/tinygo-go-path:/home/tinygo/go \
#		-w /src \
#		tinygo/tinygo:0.31.2 \
#		sh -c "chown -R 1000:1000 /home/tinygo/.cache && chown -R 1000:1000 /home/tinygo/go && tinygo build -x -o tmp/x.wasm -target=wasi ./..."

#docker run -dit --name tinygo-builder \
#  -v "$PWD":/src \
#  -v tinygo-modcache:/go \
#  -v tinygo-cache:/home/tinygo/.cache \
#  -w /src \
#  tinygo/tinygo:0.31.2 sleep infinity
#
## then:
#docker exec tinygo-builder tinygo build -o bin.wasm -target=wasi ./...


#docker run --rm -u root \
#  -v ${HOME}/.cache/myCache:/home/tinygo/.cache \
#  -w /src \
#  tinygo/tinygo:0.31.2 \
#  sh -c "chown -R 1000:1000 /home/tinygo/.cache && tinygo build -o bin.wasm -target=wasi ./main.go"


#docker run --rm \
#  -v ${HOME}/.cache/myCache:/home/tinygo/.cache \
#  -v $(pwd):/src \
#  -w /src \
#  tinygo/tinygo:0.31.2 \
#  sh -c "chown -R root:root /home/tinygo/.cache && chown -R 1000:1000 /home/tinygo/.cache/myCache && tinygo build -o bin.wasm -target=wasi ./main.go"


.PHONY: run
run: check build
	@echo "Running WASM module..."
	@cat fn-call.json | bky-as attest-fn-call >tmp/out.json
	@echo "Log:"
	@jq -r '.transitive_attested_function_call.logs | @base64d' tmp/out.json
	@echo "Output:"
	@jq -r '.transitive_attested_function_call.claims.output | @base64d ' tmp/out.json

.PHONY: clean
clean:
	@rm -rf tmp