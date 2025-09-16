package headers

import (
	"errors"
	"fmt"
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
	aTrim := strings.TrimSpace(after)
	if _, ok = h[bTrim]; !ok {
		h.Set(bTrim, aTrim)
	} else {
		h.Add(bTrim, aTrim)
	}

	return len(split) + 2, false, nil
}

func (h Headers) Get(key string) (string, bool) {
	lKey := strings.ToLower(key)
	val := ""
	ok := false
	switch lKey {
	case "a-im":
		val, ok = h["A-IM"]
	case "te":
		val, ok = h["TE"]
	case "x-csrf-token":
		val, ok = h["X-CSRF-Token"]
	default:
		splitKey := strings.Split(lKey, "-")
		for i, part := range splitKey {
			r := []rune(part)
			r[0] = unicode.ToUpper(r[0])
			splitKey[i] = string(r)
		}
		k := strings.Join(splitKey, "-")
		val, ok = h[k]
	}

	return val, ok
}

func (h Headers) Set(key, val string) {
	lKey := strings.ToLower(key)
	switch lKey {
	case "a-im":
		h["A-IM"] = val
	case "te":
		h["TE"] = val
	case "x-csrf-token":
		h["X-CSRF-Token"] = val
	default:
		splitKey := strings.Split(lKey, "-")
		for i, part := range splitKey {
			r := []rune(part)
			r[0] = unicode.ToUpper(r[0])
			splitKey[i] = string(r)
		}
		k := strings.Join(splitKey, "-")
		h[k] = val
	}
}

func (h Headers) Add(key, val string) {
	lKey := strings.ToLower(key)
	switch lKey {
	case "a-im":
		h["A-IM"] = fmt.Sprintf("%s, %s", h["A-IM"], val)
	case "te":
		h["TE"] = fmt.Sprintf("%s, %s", h["TE"], val)
	case "x-csrf-token":
		h["X-CSRF-Token"] = fmt.Sprintf("%s, %s", h["X-CSRF-Token"], val)
	default:
		splitKey := strings.Split(lKey, "-")
		for i, part := range splitKey {
			r := []rune(part)
			r[0] = unicode.ToUpper(r[0])
			splitKey[i] = string(r)
		}
		k := strings.Join(splitKey, "-")
		h[k] = fmt.Sprintf("%s, %s", h[k], val)
	}
}
