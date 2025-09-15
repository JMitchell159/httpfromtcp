package response

import (
	"fmt"
	"io"

	"github.com/JMitchell159/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	OK          StatusCode = 200
	ClientError StatusCode = 400
	ServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	out := []byte{}
	out = fmt.Appendf(out, "HTTP/1.1 %d ", statusCode)
	switch statusCode {
	case OK:
		out = fmt.Append(out, "OK\r\n")
	case ClientError:
		out = fmt.Append(out, "Bad Request\r\n")
	case ServerError:
		out = fmt.Append(out, "Internal Server Error\r\n")
	}

	_, err := w.Write(out)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf("%d", contentLen)
	h["Connection"] = "close"
	h["Content-Type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	out := []byte{}
	for key, val := range headers {
		out = fmt.Appendf(out, "%s: %s\r\n", key, val)
	}
	out = fmt.Append(out, "\r\n")
	_, err := w.Write(out)
	return err
}
