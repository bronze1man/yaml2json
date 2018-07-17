package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	goyaml "gopkg.in/yaml.v2"
)

func main() {
	err := translate(os.Stdin, os.Stdout)
	if err == io.EOF {
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(0)
}

func translate(in io.Reader, out io.Writer) error {
	decoder := goyaml.NewDecoder(in)

	for {
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			return err
		}
		err = transformData(&data)
		if err != nil {
			return err
		}

		output, err := json.Marshal(data)
		if err != nil {
			return err
		}

		data = nil
		_, err = out.Write(output)
		if err != nil {
			return err
		}

		_, err = io.WriteString(out, "\n")
		if err != nil {
			return err
		}
	}
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
