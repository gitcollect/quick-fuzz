package main

import (
	"fmt"
	"github.ihrint.com/quickio/quickigo"
	"log"
	"math/rand"
	"time"
)

func callbacks(spawn int) {
	for i := 0; i < spawn; i++ {
		go callback()
	}
}

func callbacks_fuzz_recv(spawn int) {
	for i := 0; i < spawn; i++ {
		go callback_fuzz_recv()
	}
}
func callbacks_fuzz_send(spawn int) {
	for i := 0; i < spawn; i++ {
		go callback_fuzz_send()
	}
}

func callback() {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		qio.Send("/qio/ping", nil,
			func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
				if code != quickigo.CODE_OK {
					log.Println("Ping failed:", code)
				}
				chCb <- true
			})

		<-chCb
	}
}

func callback_fuzz_recv() {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		var path string
		switch rand.Intn(3) {
		case 0:
			path = path_rand()
		case 1:
			path = path_valid()
		case 2:
			path = path_valid_with_rand()
		}

		qio.Send(path, nil,
			func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
				chCb <- true
			})

		<-chCb
	}
}

func callback_fuzz_send() {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		cbId := rand.Int63n(0xffffffffffff)
		qio.Send(fmt.Sprintf("/qio/callback/%d", cbId), cbId,
			func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
				chCb <- true
			})

		<-chCb
	}
}
