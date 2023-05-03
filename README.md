# go-jl

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
