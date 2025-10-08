package headers

import (
	"errors"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func getKeyValue(fl string) (string, string, error) {
	kvIdx := strings.Index(fl, ":")
	if kvIdx == -1 {
		return "", "", errors.New("Invalid field line")
	}
	k := fl[:kvIdx]
	k = strings.ToLower(k)
	if k[len(k)-1] == ' ' {
		// invalid key; contains whitespace after colon
		return "", "", errors.New("Invalid field line key")
	}

	f := func(r rune) bool {
		sp := "!#$%&'*+-.^_`|~"
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			if !strings.ContainsRune(sp, r) {
				return true
			}
		}
		return false
	}

	if strings.IndexFunc(k, f) != -1 {
		return "", "", errors.New("Invalid filed line key character")
	}
	v := strings.TrimSpace(fl[kvIdx+1:])
	return k, v, nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// n number of bytes consumed
	// done if its finish (when CRLF is on the start of the line)
	d := string(data)
	bc := 0
	for {
		idx := strings.Index(d[bc:], "\r\n")
		if idx == -1 {
			// not enough data
			return 0, false, nil
		} else if idx == 0 {
			// start of line; finished reading headers
			return bc, true, nil
		}

		// field line
		fl := strings.TrimSpace(d[bc : bc+idx])
		key, value, err := getKeyValue(fl)
		if err != nil {
			return 0, false, err
		}
		h[key] = value
		bc += idx + 2
	}
}
