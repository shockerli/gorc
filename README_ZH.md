# gorc

基于 [GoReplay](https://github.com/buger/goreplay) 实现的 HTTP API 流量录制与重放中间件，可用于服务接口迁移、重构、测试等，比如对比重构接口的输入输出是否一致。

中文 | [English](README.md)

## 依赖

`gorc` 是一个 `GoReplay` 中间件，所以需先安装 `GoReplay`。

- [或] [下载最新的二进制](https://github.com/buger/goreplay/releases)
- [或] [按文档自行编译二进制](https://github.com/buger/goreplay/wiki/Compilation)
- [或] 如果是 macOS 系统，直接使用 `brew install gor` 安装

## 安装

- 下载预编译的二进制

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

## 使用

> `--input-raw=":8001"`: 被捕获的原服务端口
>
> `--output-http="http://127.0.0.1:8002"`: 流量请求重放的新服务地址
>
> `--middleware="${path-of-gorc} ${command-or-script}"`: `GoReplay` 中间件命令配置，`gorc` 路径 + 命令/脚本

- PHP

[PHP](examples/script.php) 示例就是个对比两个接口响应结构是否一致的脚本。

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

- 任何其他命令、软件、或机器上所支持的语言脚本

## 数据

`gorc` 传递给脚本/命令的每个请求的数据，包括原请求、原响应、重放响应。

自己写的脚本，每次从 `STDIN` 读取的每一行数据就是如下 JSON 格式：

```json5
{
  // GoReplay 生成的唯一请求ID
  "req_id": "f33e1bab7f0000013e9b304d",
  // 整个用时，包括重放流量
  "latency": 1137963000,
  "request": {
    // 记录的时间戳（纳秒）
    "time": 1636958850332299000,
    // 耗时（纳秒）
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
    // 请求体
    "body": {
      "key": "kkk",
      "type": 1
    }
  },
  // 原始请求的响应信息
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
  // 重放请求的响应信息
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
    // 重放请求的响应体
    "body": {
      "code": 200,
      "data": [],
      "msg": "success"
    }
  }
}
```

## 协议

本项目基于 [MIT](LICENSE) 协议开放源代码。
