IMAGE_REGISTRY=ghcr.io/sks/kihocche
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "-dirty" || echo "")
GIT_VERSION=${GIT_COMMIT}${GIT_DIRTY}

.PHONY: setup

deps:
	@which -a wire > /dev/null || go install github.com/google/wire/cmd/wire@latest

setup: deps go/tv

go/tv:
	@go mod tidy
	@go mod vendor

build: setup
	@go build -o bin/kihocche .

run: build
	./bin/kihocche journey

docker: docker/build
	
docker/build:
	docker buildx build \
		--build-arg GIT_TAG=${GIT_VERSION} \
		-t ${IMAGE_REGISTRY}:${GIT_VERSION} \
		-f Dockerfile .
