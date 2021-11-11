# gorc

> [GoReplay](https://github.com/buger/goreplay) middleware for compare response, support any program languages

## Usage

> `--input-raw=":8001"`: original service port, which be recorded
> `--output-http="http://127.0.0.1:8002"`: replay request to another service

- PHP

```shell
gor --input-raw=":8001" --input-raw-track-response --output-http="http://127.0.0.1:8002" --middleware="go run gorc.go php examples/script.php" --output-http-track-response
gor --input-raw=":8001" --input-raw-track-response --output-http="http://127.0.0.1:8002" --middleware="go run gorc.go ./examples/script.php" --output-http-track-response
gor --input-raw=":8001" --input-raw-track-response --output-http="http://127.0.0.1:8002" --middleware="/path/to/bin/gorc ./examples/script.php" --output-http-track-response
```
