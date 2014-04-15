package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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
	ws := newWebSocket()
	cbId := 1

	for {
		utilPause()

		err := ws.expectPrefix(
			fmt.Sprintf("/qio/ping:%d=null", cbId),
			fmt.Sprintf("/qio/callback/%d", cbId))
		if err != nil {
			log.Println("Callback failed:", err)
		}

		cbId++
	}
}

func callbackFuzzRecv() {
	ws := newWebSocket()

	for {
		utilPause()

		err := ws.expectPrefix(
			fmt.Sprintf("%s:1=null", utilPath()),
			"/qio/callback/1")
		if err != nil {
			log.Println("Callback failed:", err)
		}
	}
}

func callbackFuzzSend() {
	ws := newWebSocket()

	for {
		utilPause()

		cbId := rand.Int63() ^ (rand.Int63() << 1)

		err := ws.expectPrefix(
			fmt.Sprintf("/qio/callback/%d:1=null", cbId),
			"/qio/callback/1")
		if err != nil {
			log.Println("Callback failed:", err)
		}
	}
}

func callbackFuzzChain() {
	ws := newWebSocket()

	for {
		utilPause()

		err := ws.expectPrefix(
			"/quick-fuzz/callback:1=null",
			"/qio/callback/1")
		if err != nil {
			log.Println("Callback failed:", err)
		} else {
			start := bytes.Index(ws.buff, []byte(":"))
			end := bytes.Index(ws.buff, []byte("="))
			if start != -1 && end != -1 {
				cbId, _ := strconv.ParseUint(string(ws.buff[start+1:end]), 10, 64)

				if rand.Intn(3) > 0 {
					ws.send(fmt.Sprintf("/qio/callback/%d:0=null", cbId))
				}
			}
		}
	}
}
