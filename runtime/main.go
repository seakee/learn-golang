package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
)

func main() {
	// 创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// 启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()

	ch := make(chan string)
	go func() {
		ch <- "EDDYCJY"
	}()

	<-ch

	// main
	fmt.Println("Hello World")
	fmt.Println(runtime.NumCPU())
}
