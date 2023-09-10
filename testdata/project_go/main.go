package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("GoVersion==%s", runtime.Version())
}
