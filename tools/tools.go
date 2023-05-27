//go:build tools
// +build tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://marcofranssen.nl/manage-go-tools-via-go-modules
// https://www.alexedwards.net/blog/using-go-run-to-manage-tool-dependencies
// https://blog.devgenius.io/sort-go-imports-acb76224dfa7
package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/incu6us/goimports-reviser/v3"
	_ "github.com/segmentio/golines"
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "mvdan.cc/gofumpt"
)
