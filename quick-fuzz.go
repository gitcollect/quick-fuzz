package main

import (
	"flag"
	"log"
	"runtime"
	"time"
)

var (
	runTime     = -1
	testPercent = 100
)

func init() {
	flag.IntVar(&runTime, "runTime", -1,
		"How long the fuzzer should run for, in seconds. -1 is infinite.")
	flag.IntVar(&testPercent, "testPercent", 100,
		"The percent of clients that should run for the fuzz.")

	if testPercent == 0 || testPercent > 100 {
		log.Panicf("testPercent must be between 1 and 100, %d is not valid", testPercent)
	}
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	percent := float32(testPercent) / 100

	callbacks(int(15000 * percent))
	callbacksFuzzChain(int(15000 * percent))
	callbacksFuzzRecv(int(15000 * percent))
	callbacksFuzzSend(int(15000 * percent))
	heartbeaters(int(20000 * percent))
	insanes(int(50000 * percent))
	rawsFuzzFormatted(int(15000 * percent))
	rawsFuzzFramed(int(15000 * percent))
	rawsFuzzRandom(int(15000 * percent))
	reconnectors(int(5000 * percent))
	subscribers(int(15000 * percent))
	subscribersFuzz(int(15000 * percent))

	if runTime != -1 {
		<-time.After(time.Second * time.Duration(runTime))
	} else {
		select {}
	}
}
