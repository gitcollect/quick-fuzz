package main

import (
	"fmt"
	"log"
	"math/rand"
)

func subscribers(spawn int) {
	for i := 0; i < spawn; i++ {
		go subscriber(spawn)
	}
}

func subscriber(spawned int) {
	for {
		utilPause()
		ws := newWebSocket()

		for i := 0; i < 10; i++ {
			utilPause()

			switch rand.Intn(2) {
			case 0:
				path := fmt.Sprintf("/quick-fuzz/delayed/%d", rand.Intn(spawned))

				err := ws.expect(
					fmt.Sprintf("/qio/on:1=\"%s\"", path),
					"/qio/callback/1:0={\"code\":200,\"data\":null}")
				if err != nil {
					log.Println("Subscribe failed:", err)
					continue
				}

				if rand.Intn(2) == 0 {
					utilPause()

					err = ws.expect(
						fmt.Sprintf("/qio/off:2=\"%s\"", path),
						"/qio/callback/2:0={\"code\":200,\"data\":null}")
					if err != nil {
						log.Println("Unsubscribe failed:", err)
					}
				}

			default:
				var path string
				if rand.Intn(2) == 0 {
					path = fmt.Sprintf("/qio/off:1=\"%s\"", utilPathRand())
				} else {
					path = fmt.Sprintf("/qio/on:1=\"%s\"", utilPathRand())
				}

				err := ws.expect(
					path,
					"/qio/callback/1:0={\"code\":404,\"data\":null,\"err_msg\":null}")
				if err != nil {
					log.Println("Didn't get 404:", err)
				}
			}
		}

		ws.close()
	}
}
