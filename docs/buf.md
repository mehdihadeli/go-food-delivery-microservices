# Buf

[Buf](https://buf.build/) is a tool for **Protobuf** files:

- [Linter](https://buf.build/docs/lint-usage) that enforces good API design choices and structure.
- [Breaking change detector](https://buf.build/docs/breaking-usage) that enforces compatibility at the source code or wire level
- Configurable file [builder](https://buf.build/docs/build-overview) that produces [Images](https://buf.build/docs/build-images) our extension of [FileDescriptorSets](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/descriptor.proto)

## Prerequisites

```bash
# buf: proto tool https://buf.build/docs/tour-1
brew tap bufbuild/buf
brew install buf
# or use `go get` to install Buf
GO111MODULE=on go get github.com/bufbuild/buf/cmd/buf
```

## Developer Workflow

### Info

```bash
# To list all files Buf is configured to use:
buf ls-files
# To see your currently configured lint or breaking checkers:
buf check ls-lint-checkers
buf check ls-breaking-checkers
# To see all available lint checkers independent of configuration/defaults:
 buf check ls-lint-checkers --all
```

### Build

```bash
# check
buf image build.sh -o /dev/null
buf image build.sh -o image.bin
```

### Lint

```bash
buf check lint
# We can also output errors in a format you can then copy into your buf.yaml file
buf check lint --error-format=config-ignore-yaml
# Run breaking change detection
buf check breaking --against-input image.bin
```

### Format

```bash
make proto_format
```

### Generate

```bash
make proto
```

## Tools

### grpcurl

```bash
# To use Buf-produced FileDescriptorSets with grpcurl on the fly:
grpcurl -protoset <(buf image build.sh -o -) ...
```

### ghz

```bash
# To use Buf-produced FileDescriptorSets with ghz on the fly:
ghz --protoset <(buf image build.sh -o -) ...
```

## Reference
1. [xmlking/micro-starter-kit](https://github.com/xmlking/micro-starter-kit)
1. [Style Guide](https://buf.build/docs/style-guide)
1. [Buf docs](https://buf.build/docs/introduction)
1. [Buf Example](https://github.com/bufbuild/buf-example/blob/master/Makefile)
1. [Buf Schema Registry](https://buf.build/docs/roadmap)
