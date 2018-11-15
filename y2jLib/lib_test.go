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

func mustTranslateStringWithOptions(inSlice string, options Options) (outSlice string){
	inBuf:=bytes.NewBufferString(inSlice)
	outBuf:=&bytes.Buffer{}
	err:= TranslateStreamWithOptions(inBuf,outBuf,options)
	if err!=nil{
		panic(err)
	}
	return outBuf.String()
}

func tryTranslateStringWithOptions(inSlice string, options Options) (outSlice string, err error){
	inBuf:=bytes.NewBufferString(inSlice)
	outBuf:=&bytes.Buffer{}
	err = TranslateStreamWithOptions(inBuf,outBuf,options)
    if err == nil {
        outSlice = outBuf.String()
    }
    return
}

func mustMatchOutput(input, output, expectedOutput string) {
    if output != expectedOutput {
        panic("Failure.\n-- Input:\n"+input+"\n-- Output:\n"+output+"\n-- Expected Output:\n"+expectedOutput+"\n")
    }
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
}
	for _,cas:=range casList{
		out:=mustTranslateStreamString(cas.in)
        mustMatchOutput(cas.in, out, cas.out)
	}
}

func TestStrictParse(ot *testing.T) {
    good_in := " one: 1\n two: 2"
    bad_in := " one: 1\n one: 2"

    options := Options{}

    // Verify that both inputs parse without strict parsing
    mustTranslateStringWithOptions(good_in, options)
    mustTranslateStringWithOptions(bad_in, options)

    options.ParseStrict = true

    // Under strict parsing, good input should succeed.
    mustTranslateStringWithOptions(good_in, options)

    // Under strict parsing, bad input should fail.
    _, err := tryTranslateStringWithOptions(bad_in, options)
    if err == nil {
        panic("Strict parsing should have failed but it didn't.")
    }
}

func TestPrettyPrint(ot *testing.T) {
    input :=
`one: 1
two: 2`

    expected_output :=
`{
  "one": 1,
  "two": 2
}
`
    options := Options { Indent: "  " }

    output := mustTranslateStringWithOptions(input, options)

    mustMatchOutput(input, output, expected_output)
}
