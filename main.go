package main

import (
	"fmt"
	"runtime/debug"

	"petezalew.ski/pit/cmd"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()
	cmd.Execute()
}
