package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	RequestPayload          = '1'
	ResponsePayload         = '2'
	ReplayedResponsePayload = '3'
)

type Request struct {
	Time   time.Duration `json:"time"`
	Method string        `json:"method"`
	URI    string        `json:"uri"`
	Proto  string        `json:"proto"`
	Body   string        `json:"body"`
}

type Response struct {
	Time       time.Duration `json:"time"`
	Header     http.Header   `json:"header"`
	Status     string        `json:"status"`
	StatusCode int           `json:"status_code"`
	Proto      string        `json:"proto"`
	Body       string        `json:"body"`
}

type Capture struct {
	Request          Request  `json:"request"`
	OriginalResponse Response `json:"original_response"`
	ReplayedResponse Response `json:"replayed_response"`
}

var store = map[string]*Capture{}

func main() {
	// scan stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		encoded := scanner.Bytes()
		buf := make([]byte, len(encoded)/2)
		hex.Decode(buf, encoded)

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

	_, ok := store[reqID]
	if !ok {
		store[reqID] = &Capture{}
	}
	pt, _ := time.ParseDuration(string(meta[2]))

	Debug("Received payload:", string(buf))

	switch payloadType {
	case RequestPayload: // Request
		req := parseReq(payload)
		req.Time = pt
		store[reqID].Request = req

		// Emitting data back
		os.Stdout.Write(encode(buf))

	case ResponsePayload: // Original response
		res := parseRes(payload)
		res.Time = pt
		store[reqID].OriginalResponse = res

	case ReplayedResponsePayload: // Replayed response
		res := parseRes(payload)
		res.Time = pt
		store[reqID].ReplayedResponse = res

		// TODO: callback
		//spew.Fdump(os.Stderr, store[reqID])

		// del data
		delete(store, reqID)
	}
}

func encode(buf []byte) []byte {
	dst := make([]byte, len(buf)*2+1)
	hex.Encode(dst, buf)
	dst[len(dst)-1] = '\n'

	return dst
}

func Debug(args ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		fmt.Fprintf(os.Stderr, "[%s] ", time.Now().Format("2006-01-02 15:04:05.000000"))
		fmt.Fprintln(os.Stderr, args...)
	}
}

func parseReq(buf []byte) Request {
	req, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	var p []byte
	if req.Body != nil {
		defer req.Body.Close()
		p, _ = ioutil.ReadAll(req.Body)
	}
	return Request{
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
		defer res.Body.Close()
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
