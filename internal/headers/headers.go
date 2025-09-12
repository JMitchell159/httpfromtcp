package headers

import (
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	h := make(Headers)
	return h
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	s := string(data)
	if !strings.Contains(s, "\r\n") {
		return 0, false, nil
	}
	if s[:2] == "\r\n" {
		return 0, true, nil
	}

	split, _, _ := strings.Cut(s, "\r\n")
	before, after, ok := strings.Cut(split, ":")
	if !ok {
		return 0, false, errors.New("there must be a colon present for each header")
	}
	bTrim := strings.TrimLeft(before, " ")
	bTrimSplit := strings.Split(bTrim, " ")
	if len(bTrimSplit) > 1 {
		return 0, false, errors.New("there must be no whitespace between the field name and the colon")
	}
	aTrim := strings.TrimSpace(after)
	r := []rune(bTrim)
	r[0] = unicode.ToUpper(r[0])
	bTrim = string(r)
	h[bTrim] = aTrim

	return len(split) + 2, false, nil
}
