LOCALBIN := ${PWD}/bin

.PHONY: smartimports
smartimports: export SMARTIMPORTS := ${LOCALBIN}/smartimports
smartimports:
	test -f ${SMARTIMPORTS} || GOBIN=${LOCALBIN} go install github.com/pav5000/smartimports/cmd/smartimports@latest
	PATH=${PATH}:${LOCALBIN} ${SMARTIMPORTS} -path . -exclude ./static ./../_mock -local github.com/aridae/go-metrics-store

.PHONY: generate-mocks
generate-mocks: export MOCKGEN := ${LOCALBIN}/mockgen
generate-mocks:
	test -f ${MOCKGEN} || GOBIN=${LOCALBIN} go install go.uber.org/mock/mockgen@latest
	PATH=${PATH}:${LOCALBIN} go generate -run mockgen $(find . -d -name '_mock')

.PHONY: lint
lint: export GOLANGCILINT := ${LOCALBIN}/golangci-lint
lint:
	test -f ${GOLANGCILINT} || GOBIN=${LOCALBIN} go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	PATH=${PATH}:${LOCALBIN} ${GOLANGCILINT} run

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test ./...

.PHONY: build-server
build-server: export SERVERBIN := ${LOCALBIN}/server
build-server:
	go build -o ${SERVERBIN} cmd/server/main.go

.PHONY: build-agent
build-agent: export AGENTBIN := ${LOCALBIN}/agent
build-agent:
	go build -o ${AGENTBIN} cmd/agent/main.go


.PHONY: autotest_14
autotest_14: export METRICSTEST := ${LOCALBIN}/metricstest
autotest_14: export AGENTBIN := ${LOCALBIN}/agent
autotest_14: export SERVERBIN := ${LOCALBIN}/server
autotest_14: build-agent build-server
	PATH=${PATH}:${LOCALBIN} ${METRICSTEST} -test.v \
	-test.run=^TestIteration14$$ \
	-agent-binary-path=${AGENTBIN} \
	-binary-path=${SERVERBIN} \
	-server-port=8080 \
	-source-path=. \
	-file-storage-path=/tmp/metrics-tests-db.json \
	-database-dsn=postgresql://metrics-store-user:pass@localhost:5432/metrics-store \
	-key=123

