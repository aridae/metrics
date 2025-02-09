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