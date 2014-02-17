package main

import (
	"fmt"
	"github.ihrint.com/quickio/quickigo"
	"log"
	"math/rand"
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
func callbacks_fuzz_chain(spawn int) {
	for i := 0; i < spawn; i++ {
		go callback_fuzz_chain()
	}
}

func callback() {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		util_pause()

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
		util_pause()

		qio.Send(util_path(), nil,
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
		util_pause()

		cbId := uint32(rand.Intn(256)<<16 | rand.Intn(0xffff))
		qio.Send(fmt.Sprintf("/qio/callback/%d", cbId), cbId,
			func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
				chCb <- true
			})

		<-chCb
	}
}

func callback_fuzz_chain() {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		util_pause()

		cbId := uint32(rand.Intn(256)<<16 | rand.Intn(0xffff))
		qio.Send(fmt.Sprintf("/qio/callback/%d", cbId), cbId,
			func(_ interface{}, cb quickigo.ServerCbFn, code int, _ string) {
				cb(nil, nil)
				chCb <- true
			})

		<-chCb
	}
}
