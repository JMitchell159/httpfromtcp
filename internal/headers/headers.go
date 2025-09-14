package headers

import (
	"errors"
	"strings"
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
		return 2, true, nil
	}

	split, _, _ := strings.Cut(s, "\r\n")
	before, after, ok := strings.Cut(split, ":")
	if !ok {
		return 0, false, errors.New("there must be a colon present for each header")
	}
	if len(before) == 0 {
		return 0, false, errors.New("the field name must have a length of at least 1")
	}
	bTrim := strings.TrimLeft(before, " ")
	if len(bTrim) == 0 {
		return 0, false, errors.New("the field name must have at least 1 non-space character directly preceeding the colon")
	}
	bTrimSplit := strings.Split(bTrim, " ")
	if len(bTrimSplit) > 1 {
		return 0, false, errors.New("there must be no whitespace between the field name and the colon")
	}
	for _, c := range bTrimSplit[0] {
		if c == '!' || (c >= '#' && c <= '\'') || c == '*' || c == '+' || c == '-' || c == '.' || (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= '^' && c <= 'z') || c == '|' || c == '~' {
			continue
		}
		return 0, false, errors.New("the field name must only contain alphanumeric characters & [!, #, $, %, &, ', *, +, -, ., ^, _, `, |, ~]")
	}
	bTrim = strings.ToLower(bTrim)
	aTrim := strings.TrimSpace(after)
	if _, ok = h[bTrim]; !ok {
		h[bTrim] = aTrim
	} else {
		slice := []string{h[bTrim], aTrim}
		h[bTrim] = strings.Join(slice, ", ")
	}

	return len(split) + 2, false, nil
}
