package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func mustCopyFile(fromPath string, toPath string) {
	content, err := ioutil.ReadFile(fromPath)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Dir(toPath), 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(toPath, content, 0777)
	if err != nil {
		panic(err)
	}
}

func mustDeleteFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
}

func mustZipDir(dir string, targetPath string) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	buf := &bytes.Buffer{}
	zipObj := zip.NewWriter(buf)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}
		w, err := zipObj.Create(relPath)
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(content)
		if err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	err = zipObj.Close()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(targetPath, buf.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
}

func mustRunCmd(cmdList []string, env map[string]string) {
	fmt.Println(">", strings.Join(cmdList, " "))
	cmd := exec.Command(cmdList[0], cmdList[1:]...)
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
