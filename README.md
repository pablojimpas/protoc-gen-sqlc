# protoc-gen-sqlc

[![Go Reference](https://pkg.go.dev/badge/github.com/pablojimpas/protoc-gen-sqlc/.svg)](https://pkg.go.dev/github.com/pablojimpas/protoc-gen-sqlc/)
[![Go CI](https://github.com/pablojimpas/protoc-gen-sqlc/actions/workflows/go.yaml/badge.svg)](https://github.com/pablojimpas/protoc-gen-sqlc/actions/workflows/go.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pablojimpas/protoc-gen-sqlc)](https://goreportcard.com/report/github.com/pablojimpas/protoc-gen-sqlc)

Protocol Buffers plugin to generate SQL schema and queries with
[sqlc](https://sqlc.dev/) annotations.

## Why?

## Install

You can download and compile the latest version from source using Go directly:

```shell
go install github.com/pablojimpas/protoc-gen-sqlc/cmd/protoc-gen-sqlc@latest
```

or you can download pre-built binaries from the [Github releases
page](https://github.com/pablojimpas/protoc-gen-sqlc/releases/latest).

## Usage

### With `buf`

To use `buf` you have to make a config file called `buf.gen.yaml` with your
options, like this:

```yaml
version: v2
plugins:
  - local: protoc-gen-sqlc
    out: out
```

And then run `buf generate`. See [the documentation on buf
generate](https://buf.build/docs/reference/cli/buf/generate#usage) for more
help.

### With `protoc`

This plugin works with `protoc` as well. Here's a basic usage example:

```shell
protoc example.proto --sqlc_out=gen
```

## General Idea

![Idea diagram](./docs/diagrams/idea.svg)

## Schema Build Process

![Generate schema sequence diagram](./docs/diagrams/schema-seq.svg)
