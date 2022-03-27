# Building Executable
# Global Setting Variable
PACKAGE_PREFIX=github.com/liyunghao/Online-Eletronic-Voting

# Artifacts
BUILD_ARTIFACTS=client server grpcConn

.PHONY: build
build: build-client build-server

.PHONY: build-client
build-client:
	go build ${PACKAGE_PREFIX}/cmd/client

.PHONY: build-server
build-server:
	go build ${PACKAGE_PREFIX}/cmd/server

.PHONY: build-test
build-test:
	go build ${PACKAGE_PREFIX}/test/grpcConn

.PHONY: compile-proto
compile-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    	   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	   internal/voting/voting.proto

.PHONY: clean
clean:
	rm -rf ${BUILD_ARTIFACTS}
