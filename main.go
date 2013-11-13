package main


import (
	"encoding/json"
	"launchpad.net/goyaml"
	"os"
	"io/ioutil"
	"fmt"
	"errors"
	"strconv"
)
func main(){
	err:=_main()
	if err!=nil{
		fmt.Fprintln(os.Stderr,err)
		os.Exit(2)
	}
	os.Exit(0)
}
func _main()error{
	var data interface{}
	input,err:=ioutil.ReadAll(os.Stdin)
	if err!=nil{
		return err
	}
	err=goyaml.Unmarshal(input,&data)
	if err!=nil{
		return err
	}
	data,err=transformData(data)
	if err!=nil{
		return err
	}

	output,err:=json.Marshal(data)
	if err!=nil{
		return err
	}
	_,err=os.Stdout.Write([]byte(output))
	return err
}
func transformData(in interface {})(out interface{},err error){
	switch in.(type){
	case map[interface {}]interface {}:
		o:=make(map[string]interface {})
		for k,v:=range in.(map[interface {}]interface {}){
			sk:=""
			switch k.(type){
			case string:
				sk=k.(string)
			case int:
			    sk=strconv.Itoa(k.(int))
			default:
				return nil,errors.New(
				fmt.Sprintf("type not match: expect map key string or int get: %T",k))
			}
			v,err= transformData(v)
			if err!=nil{
				return nil,err
			}
			o[sk] = v
		}
		return o,nil
	case []interface {}:
		in1:=in.([]interface {})
		len1:=len(in1)
		o:=make([]interface {},len1)
		for i:=0;i<len1;i++{
			o[i],err = transformData(in1[i])
			if err!=nil{
				return nil,err
			}
		}
		return o,nil
	default:
		return in,nil
	}
	return in,nil
}
