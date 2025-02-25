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

GOOS=darwin
GOARCH=amd64

VERSION = v0.0.1
COMMIT = $(shell git rev-parse HEAD)
DATE = $(shell date "+%Y-%m-%d")
LDFLAGS = -X main.buildVersion=$(VERSION) -X main.buildDate=$(DATE) -X main.buildCommit=$(COMMIT)

.PHONY: build-server
build-server: export SERVERBIN := ${LOCALBIN}/server
build-server:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ${SERVERBIN} -ldflags "$(LDFLAGS)" cmd/server/main.go

.PHONY: build-agent
build-agent: export AGENTBIN := ${LOCALBIN}/agent
build-agent:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ${AGENTBIN} -ldflags "$(LDFLAGS)" cmd/agent/main.go

.PHONY: test
test:
	go test ./...

.PHONY: test-coverage
test-coverage: export COVERAGE_OUT_FILE := ./coverage.out
test-coverage:
	go test ./... -coverpkg=./... -coverprofile=${COVERAGE_OUT_FILE} -vet=all
	go tool cover -func=${COVERAGE_OUT_FILE}

# usage: make pprof-http pprof-port=9091 app-port=8081 profile=heap seconds=30
.PHONY: pprof-http
pprof-http:
	go tool pprof -http=":${pprof-port}" -seconds=${seconds} http://127.0.0.1:${app-port}/debug/pprof/${profile}

# usage: make pprof app-port=8081 profile=heap seconds=30
.PHONY: pprof-cli
pprof-cli:
	go tool pprof -seconds=${seconds} http://127.0.0.1:${app-port}/debug/pprof/${profile}

# usage: make pprof-text app-port=8081 profile=heap seconds=30 output=base.pprof
.PHONY: pprof-text
pprof-text:
	go tool pprof --text -seconds=${seconds} http://127.0.0.1:${app-port}/debug/pprof/${profile} > ${output}


# usage: make bench out=filename.txt
.PHONY: bench
bench:
	go test ./... -bench=./... -benchmem > ${out}

# usage: make benchstat old=old.txt new=new.txt
.PHONY: benchstat
benchstat: export BENCHSTATBIN := ${LOCALBIN}/benchstat
benchstat:
	test -f ${BENCHSTATBIN} || GOBIN=${LOCALBIN} go install golang.org/x/perf/cmd/benchstat@latest
	PATH=${PATH}:${LOCALBIN} ${BENCHSTATBIN} ${old} ${new}

.PHONY: fieldalignment-fix
fieldalignment-fix: export FIELDALIGNMENTBIN := ${LOCALBIN}/fieldalignment
fieldalignment-fix:
	test -f ${FIELDALIGNMENTBIN} || GOBIN=${LOCALBIN} go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	PATH=${PATH}:${LOCALBIN} ${FIELDALIGNMENTBIN} --fix ./...

export STATICLINTBIN = ${LOCALBIN}/staticlint

.PHONY: build-staticlint
build-staticlint:
	go build -o ${STATICLINTBIN} cmd/staticlint/main.go

# usage: make staticlint pattern=./...
.PHONY: staticlint
staticlint: build-staticlint
	PATH=${PATH}:${LOCALBIN} ${STATICLINTBIN} ${pattern}

# usage: make generate-rsa-keys path=.certs/key
# output: .certs/key.pem containing private key and .certs/key.pem.pub containing public key
.PHONY: generate-rsa-keys
generate-rsa-keys:
	openssl genrsa -out ${path}.pem 4096
	openssl rsa -in ${path}.pem -outform PEM -pubout -out ${path}.pem.pub