package main

import (
	"fmt"
	"github.ihrint.com/quickio/quickigo"
	"math/rand"
	"time"
)

func subscribers(spawn int) {
	for i := 0; i < spawn; i++ {
		go subscriber(spawn)
	}
}

func subscriber(spawned int) {
	qio := quickigo.New([]string{ADDR})
	qio.Open()

	var fn quickigo.SubCb = func(data interface{}) {}

	for {
		path := fmt.Sprintf("/fuzzer/delayed/%d", rand.Intn(spawned))

		<-time.After(time.Second * time.Duration(rand.Intn(5)))
		qio.On(path, &fn)

		<-time.After(time.Second * time.Duration(rand.Intn(5)))
		qio.Off(path, nil)
	}
}
