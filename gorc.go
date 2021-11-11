package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	RequestPayload          = '1'
	ResponsePayload         = '2'
	ReplayedResponsePayload = '3'
)

type Request struct {
	StartAt time.Duration `json:"time"`
	Latency time.Duration `json:"latency"`
	Header  http.Header   `json:"header"`
	Method  string        `json:"method"`
	URI     string        `json:"uri"`
	Proto   string        `json:"proto"`
	Body    string        `json:"body"`
}

type Response struct {
	StartAt    time.Duration `json:"time"`
	Latency    time.Duration `json:"latency"`
	Header     http.Header   `json:"header"`
	Status     string        `json:"status"`
	StatusCode int           `json:"status_code"`
	Proto      string        `json:"proto"`
	Body       string        `json:"body"`
}

type Capture struct {
	ReqID            string        `json:"req_id"`
	Latency          time.Duration `json:"latency"`
	Request          Request       `json:"request"`
	OriginalResponse Response      `json:"original_response"`
	ReplayedResponse Response      `json:"replayed_response"`
}

type Script struct {
	Stdin  io.Writer
	Stdout io.Reader
}

var reqs = map[string]*Capture{}
var spt *Script

func main() {
	// exec script command
	execScript()

	// scan stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		encoded := scanner.Bytes()
		buf := make([]byte, len(encoded)/2)
		_, _ = hex.Decode(buf, encoded)

		process(buf)
	}
}

func process(buf []byte) {
	// First byte indicate payload type, possible values:
	//  1 - Request
	//  2 - Response
	//  3 - ReplayedResponse
	payloadType := buf[0]
	headerSize := bytes.IndexByte(buf, '\n') + 1
	header := buf[:headerSize-1]

	// Header contains space separated values of: request type, request id, and request start time (or round-trip time for responses)
	meta := bytes.Split(header, []byte(" "))
	// For each request you should receive 3 payloads (request, response, replayed response) with same request id
	reqID := string(meta[1])
	payload := buf[headerSize:]

	data, ok := reqs[reqID]
	if !ok {
		data = &Capture{
			ReqID: reqID,
		}
		reqs[reqID] = data
	}
	pt, _ := time.ParseDuration(string(meta[2]) + "ns")
	latency, _ := time.ParseDuration(string(meta[3]) + "ns")

	Debug("Received payload:", string(buf))

	switch payloadType {
	case RequestPayload: // Request
		req := parseReq(payload)
		req.StartAt = pt
		req.Latency = latency
		data.Request = req

		// Emitting data back
		_, _ = os.Stdout.Write(encode(buf))

	case ResponsePayload: // Original response
		res := parseRes(payload)
		res.StartAt = pt
		res.Latency = latency
		data.OriginalResponse = res

	case ReplayedResponsePayload: // Replayed response
		res := parseRes(payload)
		res.StartAt = pt
		res.Latency = latency
		data.ReplayedResponse = res
		data.Latency = data.ReplayedResponse.StartAt - data.Request.StartAt

		// script callback
		v, _ := json.Marshal(data)
		v = append(v, '\n')
		_, _ = spt.Stdin.Write(v)

		// del data
		delete(reqs, reqID)
	}
}

func execScript() {
	if len(os.Args) < 2 {
		Log("[ERROR] missing callback script")
	}

	cmd := exec.CommandContext(context.Background(), os.Args[1], os.Args[2:]...)
	cmd.Stderr = os.Stderr

	spt = new(Script)
	spt.Stdin, _ = cmd.StdinPipe()
	spt.Stdout, _ = cmd.StdoutPipe()

	go func() {
		var err error
		if err = cmd.Start(); err == nil {
			err = cmd.Wait()
		}
		if err != nil {
			if e, ok := err.(*exec.ExitError); ok {
				status := e.Sys().(syscall.WaitStatus)
				if status.Signal() == syscall.SIGKILL /*killed or context canceled */ {
					return
				}
			}
			Debug(fmt.Sprintf("[SCRIPT] command[%q] error: %q", strings.Join(os.Args[1:], " "), err.Error()))
		}
	}()
}

func encode(buf []byte) []byte {
	dst := make([]byte, len(buf)*2+1)
	hex.Encode(dst, buf)
	dst[len(dst)-1] = '\n'

	return dst
}

func Log(args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "%s ", time.Now().Format("2006-01-02 15:04:05.000000"))
	_, _ = fmt.Fprintln(os.Stderr, args...)
}

func Debug(args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		Log(args...)
	}
}

func parseReq(buf []byte) Request {
	req, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	var p []byte
	if req.Body != nil {
		defer func() { _ = req.Body.Close() }()
		p, _ = ioutil.ReadAll(req.Body)
	}
	return Request{
		Header: req.Header,
		Method: req.Method,
		URI:    req.RequestURI,
		Proto:  req.Proto,
		Body:   string(p),
	}
}

func parseRes(buf []byte) Response {
	res, _ := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
	var p []byte
	if res.Body != nil {
		defer func() { _ = res.Body.Close() }()
		p, _ = ioutil.ReadAll(res.Body)
	}
	return Response{
		Header:     res.Header,
		Status:     res.Status,
		StatusCode: res.StatusCode,
		Proto:      res.Proto,
		Body:       string(p),
	}
}
