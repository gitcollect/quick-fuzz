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

func callbacksFuzzRecv(spawn int) {
	for i := 0; i < spawn; i++ {
		go callbackFuzzRecv()
	}
}
func callbacksFuzzSend(spawn int) {
	for i := 0; i < spawn; i++ {
		go callbackFuzzSend()
	}
}
func callbacksFuzzChain(spawn int) {
	for i := 0; i < spawn; i++ {
		go callbackFuzzChain()
	}
}

func callback() {
	qio := utilCreateClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		utilPause()

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

func callbackFuzzRecv() {
	qio := utilCreateClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		utilPause()

		qio.Send(utilPath(), nil,
			func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
				chCb <- true
			})

		<-chCb
	}
}

func callbackFuzzSend() {
	qio := utilCreateClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		utilPause()

		cbId := uint32(rand.Intn(256)<<16 | rand.Intn(0xffff))
		qio.Send(fmt.Sprintf("/qio/callback/%d", cbId), cbId,
			func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
				chCb <- true
			})

		<-chCb
	}
}

func callbackFuzzChain() {
	qio := utilCreateClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		utilPause()

		cbId := uint32(rand.Intn(256)<<16 | rand.Intn(0xffff))
		qio.Send(fmt.Sprintf("/qio/callback/%d", cbId), cbId,
			func(_ interface{}, cb quickigo.ServerCbFn, code int, _ string) {
				cb(nil, nil)
				chCb <- true
			})

		<-chCb
	}
}
