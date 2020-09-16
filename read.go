package banglisp

import (
	"bufio"
	"fmt"
	"io"
)

func isSpace(c byte) bool {
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isDelimiter(c byte) bool {
	return isSpace(c) || c == '(' || c == ')' || c == '"' || c == ';'
}

func nextCharIsDigit(br *bufio.Reader) bool {
	bs, err := br.Peek(1)
	if err != nil {
		return false
	}

	return isDigit(bs[0])
}

func skipWhiteSpace(br *bufio.Reader) error {
	for {
		c, err := br.ReadByte()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if isSpace(c) {
			continue
		}

		if c == ';' { // comment
			for {
				c, err = br.ReadByte()
				if err == io.EOF {
					break
				}

				if err != nil {
					return err
				}

				if c == '\n' {
					break
				}
			}
		}

		if err := br.UnreadByte(); err != nil {
			return err
		}

		return nil
	}
}

func Read(r io.Reader) (*Object, error) {
	br := bufio.NewReader(r)

	if err := skipWhiteSpace(br); err != nil {
		return nil, err
	}

	c, err := br.ReadByte()
	if err != nil {
		return nil, err
	}

	if isDigit(c) || (c == '-' && nextCharIsDigit(br)) {
		// fixnum
		var sign int64 = 1
		if c == '-' {
			sign = -1
		} else {
			if err := br.UnreadByte(); err != nil {
				return nil, err
			}
		}

		var num int64 = 0
		eof := false
		for {
			c, err = br.ReadByte()
			if err == io.EOF {
				eof = true
				break
			}
			if err != nil {
				return nil, err
			}

			if !isDigit(c) {
				break
			}

			num = (num * 10) + int64(c-'0')
		}

		num *= sign

		if eof || isDelimiter(c) {
			if err := br.UnreadByte(); err != nil {
				return nil, err
			}

			return NewFixnum(num), nil
		}

		return nil, fmt.Errorf("could not parse fixnum")
	}

	return nil, fmt.Errorf("unsupported data type")
}
