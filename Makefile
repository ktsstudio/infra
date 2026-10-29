test: check-psql-env
	go test -short ./...

test-all: check-psql-env
	go test ./...

# update the expected command output file
test/update:
	go test ./internal/cmd -test.update-golden

DOCKER_CONTEXT=.

docker/%:
	docker buildx build $(DOCKER_CONTEXT) --load -t infrahq/$*:dev

docker/ui: DOCKER_CONTEXT=ui

dev: dev/server

dev/context:
	kubectl config use-context docker-desktop

helm/update:
	helm repo update infrahq

dev/server: dev/context docker/infra docker/ui helm/update
	helm upgrade --install --wait \
		--set-string server.image.pullPolicy=Never \
		--set-string server.image.tag=dev \
		--set-string server.podAnnotations.checksum=$$(docker images -q infrahq/infra:dev) \
		--set-string ui.image.pullPolicy=Never \
		--set-string ui.image.tag=dev \
		--set-string ui.podAnnotations.checksum=$$(docker images -q infrahq/ui:dev) \
		infra-server infrahq/infra-server \
		$(flags)

dev/connector: dev/context docker/infra
	helm upgrade --install --wait \
		--set-string connector.image.pullPolicy=Never \
		--set-string connector.image.tag=dev \
		--set-string connector.podAnnotations.checksum=$$(docker images -q infrahq/infra:dev) \
		infra infrahq/infra \
		$(flags)

dev/clean:
	kubectl config use-context docker-desktop
	helm uninstall infra-server || true
	helm uninstall infra || true

postgres:
	docker run -d --name=postgres-dev --rm \
		-e POSTGRES_PASSWORD=password123 \
		--tmpfs=/var/lib/postgresql \
		-p 127.0.0.1:15432:5432 \
		postgres:18.6-alpine@sha256:10a5c69a83035ab5239f5620d708520e8b28f5aa1f2b20143e6024532edadfb5 -c fsync=off -c full_page_writes=off \
			-c max_connections=100
	@echo
	@echo Copy the line below into the shell used to run tests
	@echo 'export POSTGRESQL_CONNECTION="host=localhost port=15432 user=postgres dbname=postgres password=password123"'


# Run security analyzers and the repository's SQL query analyzer.
GOLANGCI_LINT_VERSION ?= v2.15.0

lint:
	GOTOOLCHAIN=$$(go env GOVERSION) go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run
	(cd internal/tools/querylinter; go build -o /tmp/infra-querylinter ./cmd)
	go vet -vettool=/tmp/infra-querylinter ./...

.PHONY: docs/api/openapi3.json
docs/api/openapi3.json:
	go run -ldflags '-s -X github.com/infrahq/infra/internal.Version=0.0.0' ./internal/openapigen $@

check-psql-env:
ifndef POSTGRESQL_CONNECTION
	$(error POSTGRESQL_CONNECTION is not defined. Use `make postgres` if you need to start postgres)
endif
