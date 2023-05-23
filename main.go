package main

import (
	"os"
	"github.com/bronze1man/yaml2json/y2jLib"
)

const Version = "1.3.2"
const helpInfo = `Transform yaml string to json string without the type infomation.
Usage:
echo "a: 1" | yaml2json
yaml2json < 1.yml > 2.json
`
func main() {
	if len(os.Args)>=2{
		switch os.Args[1] {
		case "--help":
			os.Stdout.WriteString(helpInfo)
			os.Exit(1)
			return
		case "--version":
			os.Stdout.WriteString(Version+"\n")
			os.Exit(1)
			return
		default:
			os.Stderr.WriteString("not supported command line flag \n")
			os.Exit(3)
			return
		}
	}
	err := y2jLib.TranslateStream(os.Stdin, os.Stdout)
	if err == nil {
		os.Exit(0)
	}
	os.Stderr.WriteString(err.Error())
	os.Stderr.WriteString("\n")
	os.Exit(2)
}
