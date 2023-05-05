package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "code/code-user/main.go")
	var stderr, out bytes.Buffer
	cmd.Stderr = &stderr //标准错误
	cmd.Stdout = &out    //标准输出
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdinPipe.Close()
	io.WriteString(stdinPipe, "20 11\n")
	//根据测试案例运行代码，拿到结果与标准结果对比
	if err := cmd.Run(); err != nil {
		log.Fatal(err, stderr.String())
	}
	fmt.Println(out.String())
}
