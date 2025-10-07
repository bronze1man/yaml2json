package y2jLib

import (
	"bytes"
	"testing"
)

func mustTranslateStreamString(inSlice string) (outSlice string) {
	inBuf := bytes.NewBufferString(inSlice)
	outBuf := &bytes.Buffer{}
	err := TranslateStream(inBuf, outBuf)
	if err != nil {
		panic(err)
	}
	return outBuf.String()
}

func TestTranslate(ot *testing.T) {
	type tCas struct {
		in  string
		out string
	}
	var casList = []tCas{
		{`property: 0testS_1value`, `{"property":"0testS_1value"}` + "\n"},
		{"v: hi", `{"v":"hi"}` + "\n"},
		{"v: 4294967296", `{"v":4294967296}` + "\n"},
		{`a: 1
b: 2`, `{"a":1,"b":2}
`},
		{`n: bar`, `{"n":"bar"}
`},
		{`+56: bar`, `{"56":"bar"}
`},
		{`m: false`, `{"m":false}
`},
		{`m: true`, `{"m":true}
`},
		{`false: false`, `{"false":false}
`},
		{`null: null`, `{"null":null}
`},
		{`5.6: 5.6`, `{"5.6":5.6}
`},
		{`key1: .inf
key2: .nan
key3: ~
---
key4: Second
key5: Document`, `{"key1":"+Inf","key2":"NaN","key3":null}
{"key4":"Second","key5":"Document"}
`},
		{
			`x:
y:`, `{"x":null,"y":null}
`,
		}, // https://github.com/bronze1man/yaml2json/issues/15
		{
			`first:
  file: without --- at top
---
second:
  file: which is actually the last one
  next: there's nothing after the last delimiter
---`,
			`{"first":{"file":"without --- at top"}}
{"second":{"file":"which is actually the last one","next":"there's nothing after the last delimiter"}}
null
`,
		},
		{
			`properties:
  property1: FOO
  property2:`,
			`{"properties":{"property1":"FOO","property2":null}}
`,
		}, // https://github.com/bronze1man/yaml2json/issues/23
	}
	for _, cas := range casList {
		out := mustTranslateStreamString(cas.in)
		if out != cas.out {
			panic("fail in:[" + cas.in + "] thisOut:[" + string(out) + "] expect:[" + cas.out + "]")
		}
	}
}

func FuzzTranslateStreamNoPanic(f *testing.F) {
	testcases := []string{`5.6: 5.6`, `m: false`}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		inBuf := bytes.NewBufferString(orig)
		outBuf := &bytes.Buffer{}
		_ = TranslateStream(inBuf, outBuf)
	})
}
