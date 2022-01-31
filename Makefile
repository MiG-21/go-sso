GOCMD=`which go`
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

ENTRYPOINT=cmd/
BIN_NAME=go-sso
BIN_DIR=bin
BIN_OUT=${BIN_DIR}/${BIN_NAME}
GOLINT_DOCKER_TAG=v1.32.2-alpine

.PHONY: docker
all: build

# this works by creating a local vendor dir with deps in it. We don't check in
# this vendor directory, we merely pass it to the docker daemon so the build
# process inside the container doesn't have to fetch private deps (it's all
# handled in the host)
# NOTE: there are security implications if we didn't do it this way and
# instead passed the credentials into the build process. Not so much for go
# services that use a multi-stage build where none of the build stuff ends up
# in the running container, but possibly for other languages and it's just
# good practice not to send your passwords to the docker daemon
vendor:
	$(GOCMD) mod vendor

deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod download

build: .git-env
	branch=$$(grep GIT_BRANCH .git-env | awk -F '=' '{print $$2}'); \
	hash=$$(grep GIT_HASH .git-env | awk -F '=' '{print $$2}'); \
	url=$$(grep GIT_URL .git-env | awk -F '=' '{print $$2}'); \
	repo=$$(grep GIT_REPO .git-env | awk -F '=' '{print $$2}'); \
	echo "Building binaries..."; \
	CGO_ENABLED=0 GO111MODULE=on go build \
		-installsuffix "static" \
		-ldflags " \
		-X github.com/MiG-21/golang-club/app.gitBranch=$${branch} \
		-X github.com/MiG-21/golang-club/app.gitUrl=$${url} \
		-X github.com/MiG-21/golang-club/app.gitHash=$${hash} \
		-X github.com/MiG-21/golang-club/app.version=0.0.1 \
		" \
		-o ${BIN_OUT} ./cmd/...;

.git-env:
	@echo Generating .git-env ...
	@echo "GIT_URL=$$(git config remote.origin.url)" > .git-env
	@echo "GIT_BRANCH=$$(git rev-parse --abbrev-ref HEAD)" >> .git-env
	@echo "GIT_HASH=$$(git rev-parse HEAD)" >> .git-env

docker: .git-env vendor
	docker build --progress plain -f docker/Dockerfile -t golang-club-img .

docker-compose-up: .git-env vendor
	docker-compose -f docker/docker-compose.yml up -d

docker-compose-down:
	docker-compose -f docker/docker-compose.yml down -v

lint: vendor
	docker run \
		--rm \
		--volume $(PWD):/src \
		--workdir /src \
		golangci/golangci-lint:${GOLINT_DOCKER_TAG} golangci-lint --timeout=1h run $(ENTRYPOINT)

test: install_tools vendor
	@echo Starting tests...
	@APP_ENV=test go test ./...

clean:
	$(GOCLEAN) $(ENTRYPOINT)
	rm -rf $(BIN_DIR) .git-env vendor || true

install_tools:
	@for package in $$(grep '_ \"' tools/tools.go | sed 's/_ //g' | sed 's/[^a-zA-Z0-9/.]//g'); do \
		echo "Installing package $${package} or skipping if already installed..."; \
		go install $${package}; \
	done

template:
	qtc -dir=templates/

doc:
	`which swag` init --parseDependency --dir=./cmd --output=./api/docs
