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
		{`n: bar`,`{"n":"bar"}
`},
		{`+56: bar`,`{"56":"bar"}
`},
		{`m: false`,`{"m":false}
`},
		{`m: true`,`{"m":true}
`},
		{`false: false`,`{"false":false}
`},
		{`null: null`,`{"null":null}
`},
		{`5.6: 5.6`,`{"5.6":5.6}
`},
{`key1: .inf
key2: .nan
key3: ~
---
key4: Second
key5: Document`,`{"key1":".inf","key2":".nan","key3":"~"}
{"key4":"Second","key5":"Document"}
`},
}
	for _,cas:=range casList{
		out:=mustTranslateStreamString(cas.in)
		if out!=cas.out{
			panic("fail in:["+cas.in+"] thisOut:["+string(out)+"] expect:["+cas.out+"]")
		}
	}
}