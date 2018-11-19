package main

import (
	"fmt"
	"runtime"
)


func main() {
	fmt.Println("cpus:", runtime.NumCPU())
	fmt.Println("cpus:", runtime.GOARCH)
	fmt.Println("goroot:", runtime.GOROOT())
	fmt.Println("archive:", runtime.GOOS)
	pc, file, line, ok := runtime.Caller(1)
	fmt.Println(pc, file, line, ok)
}
