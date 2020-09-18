package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/syohex/banglisp"
)

func _print(bw *bufio.Writer, msg string) {
	if _, err := bw.WriteString(msg); err != nil {
		fmt.Println(err)
	}
}

func _flush(bw *bufio.Writer) {
	if err := bw.Flush(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	banglisp.Initialize()

	if len(os.Args) >= 2 {
		fmt.Printf("Usage: banglisp\n")
		return
	}

	bw := bufio.NewWriter(os.Stdout)
	for {
		_print(bw, "> ")
		_flush(bw)

		exp, err := banglisp.Read(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1) // XXX
		}

		val, err := exp.Eval()
		if err != nil {
			switch v := err.(type) {
			case *banglisp.ErrUnboundVariable:
				fmt.Printf("%s\n", v.Error())
				break
			default:
				fmt.Println(err)
				os.Exit(1) // XXX
			}

			continue
		}

		fmt.Printf("%v\n", *val)
	}
}
