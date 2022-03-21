# Fault Tolerant Computing Term Project

## Setup Dev Environment

- node `17.7.2`
- go `1.18`

```sh
# Install Golang compiler & toolchain
# Mac -> install with brew or directly download binary from the website
brew install go

# Node env can be managed with nvm package manager.
# Initialize Npm Dev Environment
npm install
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
