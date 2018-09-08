package y2jLib

import (
	"testing"
	"bytes"
)

func mustTranslateStreamString(inSlice string) (outSlice string){
	inBuf:=bytes.NewBufferString(inSlice)
	outBuf:=&bytes.Buffer{}
	err:= TranslateStream(inBuf,outBuf)
	if err!=nil{
		panic(err)
	}
	return outBuf.String()
}

func TestTranslate(ot *testing.T){
	type tCas struct{
		in string
		out string
	}
	var casList = []tCas{
		{`property: 0testS_1value`,`{"property":"0testS_1value"}`+"\n"},
		{"v: hi",`{"v":"hi"}`+"\n"},
		{"v: 4294967296",`{"v":4294967296}`+"\n"},
		{`a: 1
b: 2`,`{"a":1,"b":2}
`},
}
	for _,cas:=range casList{
		out:=mustTranslateStreamString(cas.in)
		if out!=cas.out{
			panic("fail ["+cas.in+"] ["+string(out)+"] ["+cas.out+"]")
		}
	}
}