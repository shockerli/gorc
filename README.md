# gorc

HTTP API traffic recording and replay middleware based on [GoReplay](https://github.com/buger/goreplay), can be used for
migration and refactoring testing

## Requirements

`gorc` is a middleware of `GoReplay`, so you should install `GoReplay` first.

Download the latest binary from https://github.com/buger/goreplay/releases
or [compile by yourself](https://github.com/buger/goreplay/wiki/Compilation)
or `brew install gor` in macOS

## Install `gorc`

- go install

```shell
go install github.com/shockerli/gorc@latest
```

- go build

```shell
git clone https://github.com/shockerli/gorc

cd gorc

go build .
```

## Usage

> `--input-raw=":8001"`: original service port, which be recorded
>
> `--output-http="http://127.0.0.1:8002"`: replay request to another service
>
> `--middleware="${path-of-gorc} ${command-of-script}"`

- PHP

```shell
gor \
    --input-raw-track-response \
    --output-http-track-response \
    --input-raw=":8001" \
    --output-http="http://127.0.0.1:8002" \
    --middleware="go run gorc.go php examples/script.php"

gor \
    --input-raw-track-response \
    --output-http-track-response \
    --input-raw=":8001" \
    --output-http="http://127.0.0.1:8002" \
    --middleware="go run gorc.go ./examples/script.php"

gor \
    --input-raw-track-response \
    --output-http-track-response \
    --input-raw=":8001" \
    --output-http="http://127.0.0.1:8002" \
    --middleware="/path/to/bin/gorc ./examples/script.php"
```

- NodeJS

```shell
gor \
    --input-raw-track-response \
    --output-http-track-response \
    --input-raw=":8001" \
    --output-http="http://127.0.0.1:8002" \
    --middleware="go run gorc.go node examples/script.js"

gor \
    --input-raw-track-response \
    --output-http-track-response \
    --input-raw=":8001" \
    --output-http="http://127.0.0.1:8002" \
    --middleware="go run gorc.go ./examples/script.js"

gor \
    --input-raw-track-response \
    --output-http-track-response \
    --input-raw=":8001" \
    --output-http="http://127.0.0.1:8002" \
    --middleware="/path/to/bin/gorc ./examples/script.js"
```

- Any other programming language your machine supports

## License

This project is under the terms of the [MIT](LICENSE) license.
