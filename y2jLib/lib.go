package y2jLib

import (
	"io"
	"strconv"
	"fmt"
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

type Options struct {
	ParseStrict bool
	Indent string
}

func TranslateStream(in io.Reader, out io.Writer) error {
	return TranslateStreamWithOptions(in, out, Options{})
}

func TranslateStreamWithOptions(in io.Reader, out io.Writer, options Options) error {
	decoder := yaml.NewDecoder(in)
	decoder.SetStrict(options.ParseStrict)

	for {
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			if err==io.EOF{
				return nil
			}
			return err
		}
		err = transformData(&data)
		if err != nil {
			return err
		}
		output, err := marshalToJson(data, options.Indent)
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

func marshalToJson(v interface{}, indent string) ([]byte, error) {
	if len(indent) > 0 {
		return json.MarshalIndent(v, "", indent)
	} else {
		return json.Marshal(v)
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
			case bool:
				sk = strconv.FormatBool(k.(bool))
			case nil:
				sk = "null"
			case float64:
				sk = strconv.FormatFloat(k.(float64),'f',-1,64)
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
