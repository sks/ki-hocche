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

release:
	goreleaser release --rm-dist

test:
	go test ./...

docker/tag/%: docker/build
	docker tag ${IMAGE_REGISTRY}:${GIT_VERSION} ${IMAGE_REGISTRY}:$*

helm: helm/download helm/upgrade

helm/upgrade:
	helm upgrade --install --debug --wait \
		--namespace sks \
		--values ./secrets/values.yaml \
		kihocche ./iac/helm/kihocche

helm/download:
	rm -rf ./iac
	stackgen appstack download-iac \
		--uuid ca0b9550-d787-463c-a573-8294a6f38608
	unzip -d ./iac ca0b9550-d787-463c-a573-8294a6f38608.zip
	rm -rf ca0b9550-d787-463c-a573-8294a6f38608.zip

gh/secret:
	gh secret set VALUES_YAML < ./secrets/values.yaml
