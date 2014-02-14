package main

import (
	"math/rand"
	"time"
)

func reconnectors(spawn int) {
	for i := 0; i < spawn; i++ {
		go reconnector()
	}
}

func reconnector() {
	for {
		qio := createClient()
		qio.Open()
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))
		qio.Close()
	}
}
