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

func callback() {
	ws := newWebSocket()

	for {
		utilPause()
		switch rand.Intn(2) {
		case 0:
			var cbId uint64 = uint64(rand.Int63() ^ (rand.Int63() << 1))

			err := ws.expectPrefix(
				fmt.Sprintf("%s:%d=null", utilPath(), cbId),
				fmt.Sprintf("/qio/callback/%d", cbId))
			if err != nil {
				log.Println("Callback failed:", err)
			}

		default:
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
}
