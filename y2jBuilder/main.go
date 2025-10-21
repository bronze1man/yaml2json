package main

import (
	"os"
	"path/filepath"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	mustDeleteFile(filepath.Join(pwd, "tmp", "file"))
	mustDeleteFile(filepath.Join(pwd, "tmp", "fileZip"))
	mustRunCmd([]string{"go", "test", "-v", "github.com/bronze1man/yaml2json/y2jLib"}, nil)
	for _, info := range buildInfoList {
		outputPath := filepath.Join(pwd, "bin", info.goos+"_"+info.goarch, "yaml2json")
		if info.goos == "windows" {
			outputPath += ".exe"
		}
		mustRunCmd([]string{"go", "build", "-o", outputPath, "-ldflags", "-s -w", "-gcflags=-trimpath", "github.com/bronze1man/yaml2json"}, map[string]string{
			"GOOS":   info.goos,
			"GOARCH": info.goarch,
		})
		filePath := "yaml2json_" + info.goos + "_" + info.goarch
		if info.goos == "windows" {
			filePath += ".exe"
		}
		mustCopyFile(outputPath, filepath.Join(pwd, "tmp", "file", filePath))
		filePath2 := "yaml2json"
		if info.goos == "windows" {
			filePath2 += ".exe"
		}
		mustCopyFile(outputPath, filepath.Join(pwd, "tmp", "fileZip", info.goos+"_"+info.goarch, filePath2))
	}
	mustZipDir(filepath.Join(pwd, "tmp", "fileZip"), filepath.Join(pwd, "tmp", "file", "yaml2json_all.zip"))
}

type buildInfo struct {
	goos   string
	goarch string
}

// copy from /usr/local/go/src/internal/platform/zosarch.go:10
var buildInfoList = []buildInfo{
	{"aix", "ppc64"},
	//{"android", "386"}, // android/386 requires external (cgo) linking, but cgo is not enabled
	//{"android", "amd64"},
	//{"android", "arm"},
	//{"android", "arm64"},
	{"darwin", "amd64"},
	{"darwin", "arm64"},
	{"dragonfly", "amd64"},
	{"freebsd", "386"},
	{"freebsd", "amd64"},
	{"freebsd", "arm"},
	{"freebsd", "arm64"},
	{"freebsd", "riscv64"},
	{"illumos", "amd64"},
	//{"ios", "amd64"}, // default PIE binary requires external (cgo) linking, but cgo is not enabled
	//{"ios", "arm64"},
	{"js", "wasm"},
	{"linux", "386"},
	{"linux", "amd64"},
	{"linux", "arm"},
	{"linux", "arm64"},
	{"linux", "loong64"},
	{"linux", "mips"},
	{"linux", "mips64"},
	{"linux", "mips64le"},
	{"linux", "mipsle"},
	{"linux", "ppc64"},
	{"linux", "ppc64le"},
	{"linux", "riscv64"},
	{"linux", "s390x"},
	//{"linux", "sparc64"},
	{"netbsd", "386"},
	{"netbsd", "amd64"},
	{"netbsd", "arm"},
	{"netbsd", "arm64"},
	{"openbsd", "386"},
	{"openbsd", "amd64"},
	{"openbsd", "arm"},
	{"openbsd", "arm64"},
	//{"openbsd", "mips64"}, // /usr/local/go/src/syscall/bpf_bsd.go:28:9: undefined: ioctlPtr
	{"openbsd", "ppc64"},
	{"openbsd", "riscv64"},
	{"plan9", "386"},
	{"plan9", "amd64"},
	{"plan9", "arm"},
	{"solaris", "amd64"},
	{"wasip1", "wasm"},
	{"windows", "386"},
	{"windows", "amd64"},
	{"windows", "arm"},
	{"windows", "arm64"},
}
