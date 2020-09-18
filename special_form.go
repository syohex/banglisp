package banglisp

func specialQuote(args ...*Object) *Object {
	// (quote exp)
	return args[0]
}

func specialIf(args ...*Object) *Object {
	// (if cond then else)
	cond, err := args[0].Eval()
	if err != nil {
		panic(err)
	}

	if isNull(cond) {
		if len(args) == 2 {
			return nilObj
		}

		els, err := args[2].Eval()
		if err != nil {
			panic(err)
		}

		return els
	}

	then, err := args[1].Eval()
	if err != nil {
		panic(err)
	}

	return then
}

func initSpecialForm() {
	installSpecialForm("quote", specialQuote)
	installSpecialForm("if", specialIf)
}
