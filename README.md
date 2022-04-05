# logruscheck

## Description

Linter to check if logrus logging calls are made with key-value pairs and no calls are made using the format variants like `Errorf`.

## Install

### Standalone

```bash
go install github.com/jacoelho/logruscheck
```

### [golangci-lint](https://golangci-lint.run/) plugin

Build the plugin

```bash
CGO_ENABLED=1 go build -o logruscheck.so -buildmode=plugin github.com/jacoelho/logruscheck/plugin
```

