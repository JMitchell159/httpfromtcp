package request

import (
	"errors"
	"io"
	"slices"
	"strings"
)

type state int

const (
	Initialized state = iota
	Done
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	State       state
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	currentSize := bufferSize
	buff := make([]byte, currentSize)
	readToIndex := 0
	r := &Request{
		State: Initialized,
	}

	for r.State != Done {
		if len(buff) == currentSize {
			newBuff := make([]byte, 2*currentSize)
			copy(newBuff[:currentSize], buff)
			currentSize = 2 * currentSize
			clear(buff)
			buff = newBuff
		}

		n, err := reader.Read(buff[readToIndex:])
		if errors.Is(err, io.EOF) {
			r.State = Done
			break
		}

		readToIndex += n
		n, err = r.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		newBuff := make([]byte, currentSize)
		copy(newBuff, buff[n:])
		clear(buff)
		buff = newBuff
		readToIndex -= n
	}

	return r, nil
}

func parseRequestLine(s string) (*RequestLine, int, error) {
	if !strings.Contains(s, "\r\n") {
		return nil, 0, nil
	}
	reqLine := strings.Split(s, "\r\n")[0]
	reqSplit := strings.Split(reqLine, " ")
	for i := 0; i < len(reqSplit); {
		if reqSplit[i] == "" {
			reqSplit = slices.Delete(reqSplit, i, i+1)
		} else {
			i++
		}
	}

	if len(reqSplit) != 3 {
		return nil, 0, errors.New("request line must only have 3 parts (method, request-target & version) separated by spaces")
	}

	if !slices.Contains([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}, reqSplit[0]) {
		return nil, 0, errors.New("request method is not a valid request method")
	}

	if reqSplit[2] != "HTTP/1.1" {
		return nil, 0, errors.New("only HTTP/1.1 is supported")
	}

	returnVals := RequestLine{
		HttpVersion:   "1.1",
		RequestTarget: reqSplit[1],
		Method:        reqSplit[0],
	}

	return &returnVals, len(reqLine), nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		line, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *line
		r.State = Done
		return len(data), nil
	case Done:
		return 0, errors.New("trying to read data in done state")
	default:
		return 0, errors.New("unknown state")
	}
}
