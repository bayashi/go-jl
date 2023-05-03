# go-jl

<a href="https://github.com/bayashi/go-jl/blob/main/LICENSE" title="go-jl License"><img src="https://img.shields.io/badge/LICENSE-MIT-GREEN.png" alt="MIT License"></a>
<a href="https://github.com/bayashi/go-jl/actions" title="go-jl CI"><img src="https://github.com/bayashi/go-jl/workflows/main/badge.svg" alt="go-jl CI"></a>
<a href="https://pkg.go.dev/github.com/bayashi/go-jl" title="Go go-jl package reference"><img src="https://pkg.go.dev/badge/github.com/bayashi/go-jl.svg" alt="Go Reference: go-jl"></a>

Show the "JSON within JSON" log nicely

## Usage

Simple case:

```cmd
$ echo '{"foo":"{\"bar\":\"{\\\"baz\\\":123}\"}"}' | jl
{
 "foo": {
  "bar": {
   "baz": 123
  }
 }
}
```

## Installation

    go install github.com/bayashi/go-jl/cmd/jl@latest

## License

MIT License

## Author

Dai Okabayashi: https://github.com/bayashi
