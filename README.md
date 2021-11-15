# gorc

HTTP API traffic recording and replay middleware based on [GoReplay](https://github.com/buger/goreplay), can be used for migration and refactoring testing.

English | [中文](README_ZH.md)

## Requirements

`gorc` is a middleware of `GoReplay`, so you should install `GoReplay` first.

Download the latest binary from https://github.com/buger/goreplay/releases
or [compile by yourself](https://github.com/buger/goreplay/wiki/Compilation)
or `brew install gor` in macOS

## Install `gorc`

- download binary

  https://github.com/shockerli/gorc/releases

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
> `--middleware="${path-of-gorc} ${command-or-script}"`: `gor` middleware command

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

## Datatype

Your custom script, read line from `STDIN` per request:

```json5
{
  // uuid, GoReplay generate the request unique id
  "req_id": "f33e1bab7f0000013e9b304d",
  // whole time: ns
  "latency": 1137963000,
  "request": {
    // unit: ns
    "time": 1636958850332299000,
    // unit: ns
    "latency": 0,
    "header": {
      "Accept": [
        "*/*"
      ],
      "Accept-Encoding": [
        "gzip, deflate, br"
      ],
      "Cache-Control": [
        "no-cache"
      ],
      "Connection": [
        "keep-alive"
      ],
      "Content-Length": [
        "70"
      ],
      "Content-Type": [
        "application/json"
      ],
      "User-Agent": [
        "PostmanRuntime/7.28.4"
      ]
    },
    "method": "POST",
    "uri": "/es/getTaskList",
    "proto": "HTTP/1.1",
    "body": {
      "key": "kkk",
      "type": 1
    }
  },
  "original_response": {
    "time": 1636958850404046000,
    "latency": 76000,
    "header": {
      "Connection": [
        "keep-alive"
      ],
      "Content-Type": [
        "application/json; charset=UTF-8"
      ],
      "Date": [
        "Mon, 15 Nov 2021 06:47:30 GMT"
      ],
      "Server": [
        "nginx/1.21.4"
      ],
      "X-Powered-By": [
        "PHP/7.1.33"
      ]
    },
    "status": "200 OK",
    "status_code": 200,
    "proto": "HTTP/1.1",
    "body": {
      "code": 200,
      "data": [],
      "msg": "success"
    }
  },
  "replayed_response": {
    "time": 1636958851470262000,
    "latency": 21424000,
    "header": {
      "Content-Length": [
        "107"
      ],
      "Content-Type": [
        "application/json; charset=utf-8"
      ],
      "Date": [
        "Mon, 15 Nov 2021 06:46:55 GMT"
      ]
    },
    "status": "200 OK",
    "status_code": 200,
    "proto": "HTTP/1.1",
    "body": {
      "code": 200,
      "data": [],
      "msg": "success"
    }
  }
}
```

## License

This project is under the terms of the [MIT](LICENSE) license.
