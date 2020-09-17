package banglisp

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func isSpace(c byte) bool {
	return c == ' ' || c == '\f' || c == '\n' || c == '\r' || c == '\t' || c == '\v'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDelimiter(c byte) bool {
	return isSpace(c) || c == '(' || c == ')' || c == '"' || c == ';'
}

func isInitialSymbolChar(c byte) bool {
	return isAlpha(c) || c == '*' || c == '/' || c == '>' || c == '<' || c == '=' || c == '?' || c == '!'
}

func nextCharIsDigit(br *bufio.Reader) bool {
	bs, err := br.Peek(1)
	if err != nil {
		return false
	}

	return isDigit(bs[0])
}

func nextCharIsDelimiter(br *bufio.Reader) bool {
	bs, err := br.Peek(1)
	if err != nil {
		return false
	}

	return isDelimiter(bs[0])
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

func readNumber(br *bufio.Reader, first byte) (*Object, error) {
	var sign int64 = 1
	if first == '-' {
		sign = -1
	} else {
		if err := br.UnreadByte(); err != nil {
			return nil, err
		}
	}

	var c byte
	var err error
	var num float64 = 0
	hasPoint := false
	var div float64 = 10
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

		if c == '.' {
			if hasPoint {
				return nil, fmt.Errorf("float value contains multiple dots")
			}
			hasPoint = true
			continue
		}

		if !isDigit(c) {
			break
		}

		if hasPoint {
			num = num + (float64(c-'0') / div)
			div *= 10
		} else {
			num = (num * 10) + float64(c-'0')
		}
	}

	num *= float64(sign)

	if eof || isDelimiter(c) {
		if err := br.UnreadByte(); err != nil {
			return nil, err
		}

		if hasPoint {
			return newFloat(num), nil
		} else {
			return newFixnum(int64(num)), nil
		}
	}

	return nil, fmt.Errorf("could not parse fixnum")
}

func readString(br *bufio.Reader) (*Object, error) {
	var sb strings.Builder

	for {
		c, err := br.ReadByte()
		if err == io.EOF {
			return nil, fmt.Errorf("string literal is not terminated")
		}
		if err != nil {
			return nil, err
		}

		if c == '"' {
			break
		} else if c == '\\' {
			c, err = br.ReadByte()
			if err == io.EOF {
				return nil, fmt.Errorf("string literal is not terminated")
			}
			if err != nil {
				return nil, err
			}

			if c == 'n' {
				c = '\n'
			}
		}

		sb.WriteByte(c)
	}

	return newString(sb.String()), nil
}

func readSymbol(br *bufio.Reader, c byte) (*Object, error) {
	var sb strings.Builder
	var err error
	for {
		if !(isInitialSymbolChar(c) || isDigit(c) || c == '+' || c == '-') {
			break
		}

		sb.WriteByte(c)
		c, err = br.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if !isDelimiter(c) {
		return nil, fmt.Errorf("symbol not followed by delimiter")
	}

	if err := br.UnreadByte(); err != nil {
		return nil, err
	}

	return newSymbol(sb.String()), nil
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
		return readNumber(br, c)
	} else if c == '"' {
		return readString(br)
	} else if isInitialSymbolChar(c) ||
		((c == '+' || c == '-') && nextCharIsDelimiter(br)) {
		return readSymbol(br, c)
	}

	return nil, fmt.Errorf("unsupported data type")
}
