#!/bin/bash

# In a bash script, set -e is a command that enables the "exit immediately" option. When this option is set, the script will terminate immediately if any command within the script exits with a non-zero status (indicating an error).
set -e

# `go install package@version` command works directly when we specified exact version, elsewhere it needs a `go.mod` and specifying corresponding version for each package

# https://github.com/incu6us/goimports-reviser
go install -v github.com/incu6us/goimports-reviser/v3@latest

# https://pkg.go.dev/golang.org/x/tools/cmd/goimports
go install golang.org/x/tools/cmd/goimports@latest

# https://github.com/mvdan/gofumpt
go install mvdan.cc/gofumpt@latest

# https://github.com/segmentio/golines
go install github.com/segmentio/golines@latest

# https://golangci-lint.run/usage/install/#install-from-source
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

go install google.golang.org/protobuf/proto@latest

# https://dev.to/techschoolguru/how-to-define-a-protobuf-message-and-generate-go-code-4g4e
# https://stackoverflow.com/questions/13616033/install-protocol-buffers-on-windows
go install github.com/golang/protobuf/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# https://github.com/swaggo/swag/
# https://github.com/swaggo/swag/issues/817
# swag cli v1.8.3 - upper versions have some problems with generic types
go install github.com/swaggo/swag/cmd/swag@latest
# go install github.com/swaggo/swag/cmd/swag@v1.8.3

# https://github.com/deepmap/oapi-codegen
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# https://github.com/OpenAPITools/openapi-generator-cli
npm install -g @openapitools/openapi-generator-cli

# https://vektra.github.io/mockery/latest/installation/
go install github.com/vektra/mockery/v2@latest

OS="$(uname -s)"

if [[ "$OS" == "Linux" ]]; then
    # https://github.com/bufbuild/buf
    echo "Installing Buff on Linux..."
    # Linux installation commands
    curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-Linux-x86_64 \
        -o /usr/local/bin/buf
    chmod +x /usr/local/bin/buf
    echo "Buff installed successfully."

    # https://k6.io/docs/get-started/installation/
    echo "Installing k6 on Linux..."
    sudo gpg -k
    sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
    echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
    sudo apt-get update
    sudo apt-get install k6

    # https://grpc.io/docs/protoc-installation/
    apt install -y protobuf-compiler
elif [[ "$OS" == "MINGW"* || "$OS" == "MSYS"* ]]; then
    # https://github.com/bufbuild/buf
    echo "Installing Buff on Windows..."
    # Windows installation commands
    curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-Windows-x86_64.exe \
        -o buf.exe
    Move-Item -Force buf.exe $Env:ProgramFiles\buf.exe
    echo "Buff installed successfully."

     # https://k6.io/docs/get-started/installation/
     echo "Installing k6 on Windows..."
     winget install k6

     # https://community.chocolatey.org/packages/protoc
     # https://grpc.io/docs/protoc-installation/
     choco install protoc
else
    echo "Unsupported operating system: $OS"
    exit 1
fi
