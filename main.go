package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	goyaml "gopkg.in/yaml.v2"
)

func main() {
	err := translate(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(0)
}

func translate(in io.Reader, out io.Writer) error {
	var data interface{}
	input, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	err = goyaml.Unmarshal(input, &data)
	if err != nil {
		return err
	}
	input = nil
	err = transformData(&data)
	if err != nil {
		return err
	}

	output, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = out.Write(output)
	return err
}

func transformData(pIn *interface{}) (err error) {
	switch in := (*pIn).(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			if err = transformData(&v); err != nil {
				return err
			}
			var sk string
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", k)
			}
			m[sk] = v
		}
		*pIn = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			if err = transformData(&in[i]); err != nil {
				return err
			}
		}
	}
	return nil
}
