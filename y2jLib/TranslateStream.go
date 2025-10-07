package y2jLib

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

func TranslateStream(in io.Reader, out io.Writer) error {
	decoder := yaml.NewDecoder(in)
	for {
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			if err == io.EOF {
				return nil
			}
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
	case float64:
		if math.IsInf(in, 1) {
			*pIn = "+Inf"
		} else if math.IsInf(in, -1) {
			*pIn = "-Inf"
		} else if math.IsNaN(in) {
			*pIn = "NaN"
		} else {
			// nothing to do.
		}
		return nil
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			err = transformData(&v)
			if err != nil {
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
				f := k.(float64)
				if math.IsInf(f, 1) {
					sk = "+Inf"
				} else if math.IsInf(f, -1) {
					sk = "-Inf"
				} else if math.IsNaN(f) {
					sk = "NaN"
				} else {
					sk = strconv.FormatFloat(k.(float64), 'f', -1, 64)
				}
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", k)
			}
			m[sk] = v
		}
		*pIn = m
	case map[string]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			err = transformData(&v)
			if err != nil {
				return err
			}
			m[k] = v
		}
		*pIn = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			err = transformData(&in[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
