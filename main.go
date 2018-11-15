package main

import (
	"os"
	"flag"
	"fmt"

	"github.com/bronze1man/yaml2json/y2jLib"
)

func main() {

	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "%s: Transform yaml string to json string without the type infomation.\n", os.Args[0])
		fmt.Fprintln(out, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(out,`
Examples:
	echo "a: 1" | yaml2json
	yaml2json < 1.yml > 2.json`)
	}

	parseStrict := flag.Bool("strict", false, "Enable strict YAML parsing.")
	indent := flag.String("indent", "", "If provided, pretty-prints JSON using the argument as the indent string.")

	flag.Parse()

	options := y2jLib.Options {
		ParseStrict: *parseStrict,
		Indent: *indent,
	}

	err := y2jLib.TranslateStreamWithOptions(os.Stdin, os.Stdout, options)
	if err == nil{
		os.Exit(0)
	}
	os.Stderr.WriteString(err.Error())
	os.Stderr.WriteString("\n")
	os.Exit(2)
}

