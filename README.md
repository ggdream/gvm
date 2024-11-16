# GVM is a tool for managing multiple Go versions

> Support platforms: ALL

## Installation

```bash
go install github.com/ggdream/gvm@latest
```

## Usage

1. Install a Go version

   ```bash
   gvm install 1.23.2
   ```

2. Set the Go global version

   ```bash
   gvm global 1.23.2
   ```

## Environment Variables

1. GVM_ROOT

   Store the root directory of gvm(including versions)

2. GVM_HOST

   Go software download source host. Default is <https://golang.google.cn>, you can change it to <https://go.dev> if you want to use the official source.
