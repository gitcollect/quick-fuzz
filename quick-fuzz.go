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
	totalClients = 0
	fuzzes       = [...]fuzz{
		fuzz{
			fn:  callbacks,
			num: 10000,
		},
		fuzz{
			fn:  wsHeartbeaters,
			num: 7500,
		},
		fuzz{
			fn:  httpHeartbeaters,
			num: 7500,
		},
		fuzz{
			fn:  insanes,
			num: 8000,
		},
		fuzz{
			fn:  rawsFuzzFormatted,
			num: 11000,
		},
		fuzz{
			fn:  rawsFuzzFramed,
			num: 11000,
		},
		fuzz{
			fn:  rawsFuzzRandom,
			num: 11000,
		},
		fuzz{
			fn:  reconnectors,
			num: 5000,
		},
		fuzz{
			fn:  subscribers,
			num: 11000,
		},
		fuzz{
			fn:  httpReconnectors,
			num: 11000,
		},
		fuzz{
			fn:  httpFuzzes,
			num: 11000,
		},
		fuzz{
			fn:  httpMultiRaces,
			num: 4000,
		},
	}
)

func init() {
	flag.IntVar(&runTime, "runTime", -1,
		"How long the fuzzer should run, in seconds. -1 is infinite.")
	flag.IntVar(&totalClients, "totalClients", 0,
		"Total number of clients that should be spawned. 0 is default.")
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
