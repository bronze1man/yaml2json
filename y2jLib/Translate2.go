package y2jLib

import "io"

type Translate2Req struct {
	In io.Reader
	Out io.Writer
	MultiDocumentAsJsonArray bool // https://github.com/bronze1man/yaml2json/issues/19

}
