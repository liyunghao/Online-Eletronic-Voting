# Fault Tolerant Computing Term Project

## Setup Dev Environment

- node `17.7.2`
- go `1.18`
- protoc `3.6.1`

### OSX
```sh
# Install Golang compiler & toolchain
# Mac -> install with brew or directly download binary from the website
brew install go

# Install golangci-lint for linter
brew install golangci-lint

# Install protobuf compiler
brew install protobuf

# Node env can be managed with nvm package manager.
# Initialize Npm Dev Environment
npm install
```

### Linux
```sh
sudo apt-get update

# Don't forget to check the version 
sudo apt-get install golang-go

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.45.2

# Install protobuf compiler
sudo apt-get install protobuf-compiler

# Install plugins for protobuf compiler to generate go
go get google.golang.org/protobuotoc-gen-go@v1.26
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

## Build and Run

```sh
# This will build the project and output the binary
make build

# Clean the artifacts
make clean
```

### References

Development Documents

- [golang-standards/project-layout: Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go | gRPC](https://grpc.io/docs/languages/go/)

Style Guild:

- [Git Commit Message 這樣寫會更好，替專案引入規範與範例](https://wadehuanglearning.blogspot.com/2019/05/commit-commit-commit-why-what-commit.html)
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification)
- [conventional-changelog/commitlint: 📓 Lint commit messages](https://github.com/conventional-changelog/commitlint)
- [typicode/husky: Git hooks made easy 🐶 woof!](https://github.com/typicode/husky)
- [golangci/golangci-lint: Fast linters Runner for Go](https://github.com/golangci/golangci-lint)
