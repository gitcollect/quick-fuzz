package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"
)

type fuzz struct {
	fn  func(int)
	num int
}

var (
	runTime      = -1
	totalClients = 100
	fuzzes       = [...]fuzz{
		fuzz{
			fn:  callbacks,
			num: 10000,
		},
		fuzz{
			fn:  callbacksFuzzChain,
			num: 10000,
		},
		fuzz{
			fn:  callbacksFuzzRecv,
			num: 7500,
		},
		fuzz{
			fn:  callbacksFuzzSend,
			num: 10000,
		},
		fuzz{
			fn:  heartbeaters,
			num: 7500,
		},
		fuzz{
			fn:  insanes,
			num: 6000,
		},
		fuzz{
			fn:  rawsFuzzFormatted,
			num: 10000,
		},
		fuzz{
			fn:  rawsFuzzFramed,
			num: 10000,
		},
		fuzz{
			fn:  rawsFuzzRandom,
			num: 10000,
		},
		fuzz{
			fn:  reconnectors,
			num: 5000,
		},
		fuzz{
			fn:  subscribers,
			num: 10000,
		},
		fuzz{
			fn:  subscribersFuzz,
			num: 10000,
		},
		fuzz{
			fn:  httpReconnectors,
			num: 10000,
		},
		fuzz{
			fn:  httpHeartbeaters,
			num: 7500,
		},
		fuzz{
			fn:  httpFuzzes,
			num: 10000,
		},
		fuzz{
			fn:  httpMultiRaces,
			num: 3000,
		},
	}
)

func init() {
	flag.IntVar(&runTime, "runTime", -1,
		"How long the fuzzer should run for, in seconds. -1 is infinite.")
	flag.IntVar(&totalClients, "totalClients", 0,
		"The number of clients that should be used to point everything. 0 is default.")
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	percent := 1.0

	if totalClients > 0 {
		total := 0
		for _, f := range fuzzes {
			total += f.num
		}

		percent = float64(totalClients) / float64(total)
	}

	for _, f := range fuzzes {
		clients := int(math.Ceil(float64(f.num) * percent))
		f.fn(clients)
	}

	if runTime >= 0 {
		<-time.After(time.Second * time.Duration(runTime))
	} else {
		select {}
	}
}
