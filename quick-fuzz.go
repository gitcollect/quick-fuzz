package main

import (
	"runtime"
)

const (
	ADDR = "ws://localhost:8080"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	heartbeaters(10000)
	subscribers(10000)

	select {}
}
