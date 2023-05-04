# go-jl

<a href="https://github.com/bayashi/go-jl/blob/main/LICENSE" title="go-jl License"><img src="https://img.shields.io/badge/LICENSE-MIT-GREEN.png" alt="MIT License"></a>
<a href="https://github.com/bayashi/go-jl/actions" title="go-jl CI"><img src="https://github.com/bayashi/go-jl/workflows/main/badge.svg" alt="go-jl CI"></a>
<a href="https://pkg.go.dev/github.com/bayashi/go-jl" title="Go go-jl package reference" target="_blank"><img src="https://pkg.go.dev/badge/github.com/bayashi/go-jl.svg" alt="Go Reference: go-jl"></a>

Show the "JSON within JSON" log nicely

## Usage

Simple case:

```cmd
$ echo '{"foo":"{\"bar\":\"{\\\"baz\\\":123}\"}"}' | jl -p
{
 "foo": {
  "bar": {
   "baz": 123
  }
 }
}
```

Full options:

```
Options:
  -h, --help         Display help (This message) and exit
  -p, --prettify     Prettify the JSON
  -e, --show-error   Set this option to show errors
  -v, --version      Display version and build info and exit
```

## Installation

### homebrew install

If you are using Mac:

    brew tap bayashi/go-jl
    brew install bayashi/go-jl/go-jl

### binary install

Download binary from here: https://github.com/bayashi/go-jl/releases

### go install

If you have golang envvironment:

    go install github.com/bayashi/go-jl/cmd/jl@latest

## License

MIT License

## Author

Dai Okabayashi: https://github.com/bayashi
