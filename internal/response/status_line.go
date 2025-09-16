package response

import "fmt"

type StatusCode int

const (
	OK          StatusCode = 200
	ClientError StatusCode = 400
	ServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case OK:
		reasonPhrase = "OK"
	case ClientError:
		reasonPhrase = "Bad Request"
	case ServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}
