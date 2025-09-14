package request

import (
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/JMitchell159/httpfromtcp/internal/headers"
)

type state int

const (
	Initialized state = iota
	ParsingHeaders
	ParsingBody
	Done
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
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
		State:   Initialized,
		Headers: headers.NewHeaders(),
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
			val, ok := r.Headers.Get("Content-Length")
			if ok {
				split := strings.Split(val, ", ")
				length, _ := strconv.Atoi(split[len(split)-1])
				if r.Body == nil {
					return nil, errors.New("no body, but Content-Length is present")
				}
				if length > len(r.Body) {
					return nil, errors.New("body length cannot be less than Content-Length")
				}
			}
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

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	s := string(data)
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

	return &returnVals, len(reqLine) + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != Done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}

		totalBytesParsed += n
		if totalBytesParsed == len(data) {
			return totalBytesParsed, nil
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		line, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *line
		r.State = ParsingHeaders
		return n, nil
	case ParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		if !done {
			return n, nil
		}

		r.State = ParsingBody
		if _, ok := r.Headers.Get("Content-Length"); !ok {
			r.State = Done
		}
		return n, nil
	case ParsingBody:
		sLen, _ := r.Headers.Get("Content-Length")
		split := strings.Split(sLen, ", ")
		length, err := strconv.Atoi(split[len(split)-1])
		if err != nil {
			return 0, err
		}
		if len(r.Body)+len(data) > length {
			return length, errors.New("body length cannot be greater than Content-Length")
		}
		r.Body = append(r.Body, data...)
		return len(data), nil
	case Done:
		return 0, errors.New("trying to read data in done state")
	default:
		return 0, errors.New("unknown state")
	}
}
