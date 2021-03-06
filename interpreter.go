package banglisp

import "math"

var defaultPackage *Object
var tObj *Object
var nilObj *Object
var emptyList *Object
var defaultEnvironment *Environment

func init() {
	nilObj = newSymbolInternal("nil")
	v := nilObj.value.(*Symbol)

	emptyList = newConsCell(nilObj, nilObj)

	v.value = nilObj
	v.plist = nilObj

	defaultPackage = newPackage("CL-USER")
	p := defaultPackage.value.(*Package)

	p.setSymbol(nilObj)

	v.package_ = defaultPackage

	tObj = newSymbol("t")
	tv := tObj.value.(*Symbol)
	tv.value = tObj

	piObj := newSymbol("pi")
	pv := piObj.value.(*Symbol)
	pv.value = newFloat(math.Pi)

	defaultEnvironment = newEmptyEnvironment()

	initSpecialForm()
	initBuiltinFunctions()
	initNumberFunctions()
}

func CurrentPackage() *Object {
	return defaultPackage
}

func Eval(obj *Object) (*Object, error) {
	return obj.Eval(defaultEnvironment)
}
