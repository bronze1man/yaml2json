package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"
)

func TestCmdHelp(ot *testing.T) {
	buildCmdOnce()
	cmd := exec.Command("bin/yaml2json", "--help")
	outBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	err := cmd.Run()
	ok(err != nil)
	ok(err.Error() == "exit status 1")
	ok(outBuf.String() == helpInfo, outBuf.String())
}

func TestCmdVersion(ot *testing.T) {
	// fix https://github.com/bronze1man/yaml2json/issues/16
	buildCmdOnce()
	cmd := exec.Command("bin/yaml2json", "--version")
	outBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	err := cmd.Run()
	ok(err != nil)
	ok(err.Error() == "exit status 1")
	ok(outBuf.String() == Version+"\n", outBuf.String())
}

func TestCmdExampleInReadme(ot *testing.T) {
	buildCmdOnce()
	cmd := exec.Command("bin/yaml2json")
	cmd.Stdin = bytes.NewReader([]byte("a: 1"))
	outBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	ok(outBuf.String() == `{"a":1}`+"\n", outBuf.String())
}

func TestCmdErrHandle(ot *testing.T) {
	buildCmdOnce()
	cmd := exec.Command("bin/yaml2json")
	cmd.Stdin = bytes.NewReader([]byte(`a :::: asfe234""""
sadfq23
`))
	errBuf := &bytes.Buffer{}
	cmd.Stderr = errBuf
	outBuf := &bytes.Buffer{}
	cmd.Stdout = outBuf
	err := cmd.Run()
	ok(err != nil)
	ok(err.Error() == "exit status 2")
	ok(outBuf.Len() == 0)
	ok(errBuf.String() == `yaml: line 2: could not find expected ':'`+"\n", errBuf.String())
}

var gBuildOnce sync.Once

func buildCmdOnce() {
	gBuildOnce.Do(func() {
		cmd := exec.Command("go", "build", "-o", "bin/yaml2json", "./")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	})
}

func ok(b bool, objList ...interface{}) {
	if b == false {
		panic(`fail ` + fmt.Sprintln(objList...))
	}
}
