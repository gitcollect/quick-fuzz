package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	callbacks(15000)
	callbacksFuzzChain(15000)
	callbacksFuzzRecv(15000)
	callbacksFuzzSend(15000)
	heartbeaters(20000)
	insanes(50000)
	rawsFuzzFormatted(15000)
	rawsFuzzFramed(15000)
	rawsFuzzRandom(15000)
	reconnectors(2000)
	subscribers(15000)
	subscribersFuzz(15000)

	select {}
}
