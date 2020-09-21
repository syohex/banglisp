package banglisp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
	return isAlpha(c) || c == '+' || c == '-' || c == '*' || c == '/' || c == '%' ||
		c == '>' || c == '<' || c == '=' || c == '?' || c == '!'
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

func skipWhiteSpace(br *bufio.Reader) {
	for {
		c, err := br.ReadByte()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalln(err)
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
					log.Fatalln(err)
				}

				if c == '\n' {
					break
				}
			}
		}

		unreadChar(br)
		return
	}
}

func unreadChar(br *bufio.Reader) {
	if err := br.UnreadByte(); err != nil {
		log.Fatalln(err)
	}
}

func readNumber(br *bufio.Reader, first byte) (*Object, error) {
	var sign int64 = 1
	if first == '-' {
		sign = -1
	} else {
		unreadChar(br)
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
		if !eof {
			unreadChar(br)
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
		if !(isInitialSymbolChar(c) || isDigit(c)) {
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

	unreadChar(br)

	return newSymbol(sb.String()), nil
}

func readList(br *bufio.Reader) (*Object, error) {
	skipWhiteSpace(br)

	c, err := br.ReadByte()
	if err == io.EOF {
		return nil, fmt.Errorf("list is not closed")
	}
	if err != nil {
		return nil, err
	}

	if c == ')' {
		return emptyList, nil
	}

	unreadChar(br)
	carObj, err := read1(br)
	if err != nil {
		return nil, err
	}

	skipWhiteSpace(br)

	c, err = br.ReadByte()
	if err != nil {
		return nil, err
	}

	if c == '.' {
		// dotted-pair
		if !nextCharIsDelimiter(br) {
			return nil, fmt.Errorf("dot not followed by delimiter")
		}

		cdrObj, err := read1(br)
		if err != nil {
			return nil, err
		}

		skipWhiteSpace(br)
		c, err = br.ReadByte()
		if err != nil {
			return nil, err
		}

		if c != ')' {
			return nil, fmt.Errorf("list is not closed by right paren")
		}

		return cons(carObj, cdrObj), nil
	}

	unreadChar(br)

	cdrObj, err := readList(br)
	if err != nil {
		return nil, err
	}

	return cons(carObj, cdrObj), nil
}

func read1(br *bufio.Reader) (*Object, error) {
	skipWhiteSpace(br)

	c, err := br.ReadByte()
	if err != nil {
		return nil, err
	}

	if isDigit(c) || (c == '-' && nextCharIsDigit(br)) {
		return readNumber(br, c)
	} else if c == '"' {
		return readString(br)
	} else if isInitialSymbolChar(c) {
		return readSymbol(br, c)
	} else if c == '(' {
		return readList(br)
	} else if c == '\'' {
		rest, err := read1(br)
		if err != nil {
			return nil, err
		}

		quote := intern(newString("quote"), nil)
		return cons(quote, cons(rest, emptyList)), nil
	}

	return nil, fmt.Errorf("unsupported data type")
}

func Read(r io.Reader) (*Object, error) {
	return read1(bufio.NewReader(r))
}

func ReadEvalFile(file string) (*Object, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	br := bufio.NewReader(f)
	var obj *Object
	for {
		obj, err = read1(br)
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		_, err = Eval(obj)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}
