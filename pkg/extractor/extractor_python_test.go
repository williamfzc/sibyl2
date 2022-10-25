package extractor

import (
	"testing"

	"github.com/williamfzc/sibyl2/pkg/core"
)

var pythonCode = `
import requests

def a():
	b("abcde")

@DDDDeco
@DDABC2
@CCC(abcde='acde')
def b(s):
	print("defabc")

class C(object):
	pass
`

func TestPythonExtractor_ExtractFunctions(t *testing.T) {
	parser := core.NewParser(core.LangPython)
	units, err := parser.Parse([]byte(pythonCode))
	if err != nil {
		panic(err)
	}

	extractor := GetExtractor(core.LangPython)
	functions, err := extractor.ExtractFunctions(units)
	if err != nil {
		panic(err)
	}
	for _, each := range functions {
		core.Log.Debugf("%s", each.Name)
		core.Log.Debugf("%+v", each.Extras)
	}
}
