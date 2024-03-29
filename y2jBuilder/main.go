package main

import (
	"os/exec"
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"io/ioutil"
	"archive/zip"
	"bytes"
)

func main(){
	gopath:=os.Getenv("GOPATH")
	mustDeleteFile(filepath.Join(gopath,"tmp","file"))
	mustDeleteFile(filepath.Join(gopath,"tmp","fileZip"))
	mustRunCmd([]string{"go","test","-v","github.com/bronze1man/yaml2json/y2jLib"},nil)
	type buildInfo struct{
		goos string
		goarch string
	}
	buildInfoList := []buildInfo{}
	for _,line:= range strings.Split(osAndArchList,"\n"){
		line = strings.TrimSpace(line)
		if line==``{
			continue
		}
		part:=strings.Split(line,"/")
		buildInfoList = append(buildInfoList,buildInfo{goos: part[0],goarch: part[1]})
	}
	//buildInfoList := []buildInfo{
	//	{"darwin","amd64"},
	//	{"linux","386"},
	//	{"linux","amd64"},
	//	{"linux","arm"},
	//	{"windows","386"},
	//	{"windows","amd64"},
	//	{"freebsd","386"},
	//	{"freebsd","amd64"},
	//	{"freebsd","arm"},
	//	{"netbsd","386"},
	//	{"netbsd","amd64"},
	//	{"netbsd","arm"},
	//	{"openbsd","386"},
	//	{"openbsd","amd64"},
	//	{"plan9","386"},
	//	{"linux","ppc64le"},
	//}
	for _,info:=range buildInfoList{
		outputPath := filepath.Join(gopath,"bin",info.goos+"_"+info.goarch,"yaml2json")
		if info.goos=="windows"{
			outputPath+=".exe"
		}
		mustRunCmd([]string{"go","build","-o",outputPath,"-ldflags","-s -w","-gcflags=-trimpath","github.com/bronze1man/yaml2json"},map[string]string{
			"GOOS":info.goos,
		})
		filePath:="yaml2json_"+info.goos+"_"+info.goarch
		if info.goos=="windows"{
			filePath+=".exe"
		}
		mustCopyFile(outputPath,filepath.Join(gopath,"tmp","file",filePath))
		filePath2:="yaml2json"
		if info.goos=="windows"{
			filePath2+=".exe"
		}
		mustCopyFile(outputPath,filepath.Join(gopath,"tmp","fileZip",info.goos+"_"+info.goarch,filePath2))
	}
	mustZipDir(filepath.Join(gopath,"tmp","fileZip"),filepath.Join(gopath,"tmp","file","yaml2json_all.zip"))
}

const osAndArchList = `darwin/amd64
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
illumos/amd64
linux/386
linux/amd64
linux/arm
linux/arm64
linux/mips
linux/mips64
linux/mips64l
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
openbsd/mips6
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
windows/386
windows/amd64
windows/arm
windows/arm64`
func mustCopyFile(fromPath string,toPath string){
	content,err:=ioutil.ReadFile(fromPath)
	if err!=nil{
		panic(err)
	}
	err=os.MkdirAll(filepath.Dir(toPath),0777)
	if err!=nil{
		panic(err)
	}
	err=ioutil.WriteFile(toPath,content,0777)
	if err!=nil{
		panic(err)
	}
}

func mustDeleteFile(path string){
	err:=os.RemoveAll(path)
	if err!=nil{
		if os.IsNotExist(err){
			return
		}
		panic(err)
	}
}

func mustZipDir(dir string,targetPath string){
	dir,err := filepath.Abs(dir)
	if err!=nil{
		panic(err)
	}
	buf:=&bytes.Buffer{}
	zipObj:=zip.NewWriter(buf)
	err=filepath.Walk(dir,func(path string, info os.FileInfo, err error) error{
		if err!=nil{
			panic(err)
		}
		if info.IsDir(){
			return nil
		}
		relPath,err:=filepath.Rel(dir,path)
		if err!=nil{
			panic(err)
		}
		w,err:=zipObj.Create(relPath)
		if err!=nil{
			panic(err)
		}
		content,err:=ioutil.ReadFile(path)
		if err!=nil{
			panic(err)
		}
		_,err=w.Write(content)
		if err!=nil{
			panic(err)
		}
		return nil
	})
	if err!=nil{
		panic(err)
	}
	err = zipObj.Close()
	if err!=nil{
		panic(err)
	}
	err=ioutil.WriteFile(targetPath,buf.Bytes(),0777)
	if err!=nil{
		panic(err)
	}
}

func mustRunCmd(cmdList []string,env map[string]string) {
	fmt.Println(">",strings.Join(cmdList," "))
	cmd:=exec.Command(cmdList[0],cmdList[1:]...)
	cmd.Env = os.Environ()
	for k,v:=range env{
		cmd.Env = append(cmd.Env,k+"="+v)
	}
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err:=cmd.Run()
	if err!=nil{
		panic(err)
	}
}