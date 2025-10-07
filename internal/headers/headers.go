package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
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
		kvIdx := strings.Index(fl, ":")
		if kvIdx == -1 {
			return 0, false, errors.New("Invalid field line")
		}
		k := fl[:kvIdx]
		if k[len(k)-1] == ' ' {
			// invalid key; contains whitespace after colon
			return 0, false, errors.New("Invalid field line key")
		}
		v := strings.TrimSpace(fl[kvIdx+1:])
		h[k] = v
		bc += idx + 2
	}
}
