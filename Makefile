LOCALBIN = $(PWD)/bin

.PHONY: smartimports
smartimports: export SMARTIMPORTS := ${LOCALBIN}/smartimports
smartimports:
	test -f ${SMARTIMPORTS} || GOBIN=${LOCALBIN} go install github.com/pav5000/smartimports/cmd/smartimports@latest
	PATH=${PATH}:${LOCALBIN} ${SMARTIMPORTS} -path . -exclude ./static ./../_mock -local github.com/aridae/go-metrics-store

.PHONY: generate-mocks
generate-mocks: export MOCKGEN := ${LOCALBIN}/mockgen
generate-mocks:
	test -f ${MOCKGEN} || GOBIN=${LOCALBIN} go install go.uber.org/mock/mockgen@latest
	PATH=${PATH}:${LOCALBIN} go generate -run mockgen $(shell find . -d -name '_mock')

.PHONY: lint
lint: export GOLANGCILINT := ${LOCALBIN}/golangci-lint
lint:
	test -f ${GOLANGCILINT} || GOBIN=${LOCALBIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	PATH=${PATH}:${LOCALBIN} ${GOLANGCILINT} run

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build-server
build-server: export SERVERBIN := ${LOCALBIN}/server
build-server:
	go build -o ${SERVERBIN} cmd/server/main.go

.PHONY: build-agent
build-agent: export AGENTBIN := ${LOCALBIN}/agent
build-agent:
	go build -o ${AGENTBIN} cmd/agent/main.go

.PHONY: test
test:
	go test ./...

.PHONY: test-coverage
test-coverage: export COVERAGE_OUT_FILE := ./coverage.out
test-coverage:
	go test ./... -coverpkg=./... -coverprofile=${COVERAGE_OUT_FILE} -vet=all
	go tool cover -func=${COVERAGE_OUT_FILE}

.PHONY: bench
bench:
	go test ./... -bench=./... -benchmem > ${out}

.PHONY: benchcmp
benchcmp:
	benchcmp ${old} ${new}