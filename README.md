# go-jl

<a href="https://github.com/bayashi/go-jl/blob/main/LICENSE" title="go-jl License"><img src="https://img.shields.io/badge/LICENSE-MIT-GREEN.png" alt="MIT License"></a>
<a href="https://github.com/bayashi/go-jl/actions" title="go-jl CI"><img src="https://github.com/bayashi/go-jl/workflows/main/badge.svg" alt="go-jl CI"></a>
<a href="https://goreportcard.com/report/github.com/bayashi/go-jl" title="go-jl report card"><img src="https://goreportcard.com/badge/github.com/bayashi/go-jl" alt="go-jl report card"></a>
<a href="https://pkg.go.dev/github.com/bayashi/go-jl" title="Go go-jl package reference" target="_blank"><img src="https://pkg.go.dev/badge/github.com/bayashi/go-jl.svg" alt="Go Reference: go-jl"></a>

Show the "JSON within JSON" log nicely

## Usage

`jl` command recursively converts JSON within JSON into one JSON structure.

Simple case:

```
$ cat simple.json
{
    "foo": "{\"bar\":\"{\\\"baz\\\":123}\"}"
}

$ cat simple.json | jl
{
 "foo": {
  "bar": {
   "baz": 123
  }
 }
}
```

Most use cases:

```
$ cat log.json
{
    "message": "{\"level\":\"info\",\"ts\":1557004280.5372975,\"caller\":\"zap/server_interceptors.go:40\",\"msg\":\"finished unary call with code OK\",\"grpc.start_time\":\"2019-05-04T21:11:20Z\",\"system\":\"grpc\",\"span.kind\":\"server\",\"grpc.service\":\"FooService\",\"grpc.method\":\"GetBar\",\"grpc.code\":\"OK\",\"grpc.time_ms\":248.45199584960938}\n",
    "namespace": "foo-service",
    "podName": "foo-86495899d8-m2vfl",
    "containerName": "foo-service"
}

$ cat log.json | jl
{
    "containerName": "foo-service",
    "message": {
        "caller": "zap/server_interceptors.go:40",
        "grpc.code": "OK",
        "grpc.method": "GetBar",
        "grpc.service": "FooService",
        "grpc.start_time": "2019-05-04T21:11:20Z",
        "grpc.time_ms": 248.45199584960938,
        "level": "info",
        "msg": "finished unary call with code OK",
        "span.kind": "server",
        "system": "grpc",
        "ts": 1557004280.5372975
    },
    "namespace": "foo-service",
    "podName": "foo-86495899d8-m2vfl"
}

```

Full options:

```
Options:
  -h, --help         Display help (This message) and exit
  -P, --no-prettify  Not prettify the JSON. Prettified by default
  -e, --show-error   Set this option to show errors, muted by default
  -n, --split-lf     Split line-feed \n in each element
  -t, --split-tab    Split tabs in each element
  -v, --version      Display version and build info and exit
```

## Installation

### homebrew install

If you are using Mac:

```cmd
brew tap bayashi/tap
brew install bayashi/tap/go-jl
```

### binary install

Download binary from here: https://github.com/bayashi/go-jl/releases

### go install

If you have golang envvironment:

```cmd
go install github.com/bayashi/go-jl/cmd/jl@latest
```

## License

MIT License

## Author

Dai Okabayashi: https://github.com/bayashi
