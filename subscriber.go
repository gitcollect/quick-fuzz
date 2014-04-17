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

func subscribersFuzz(spawn int) {
	for i := 0; i < spawn; i++ {
		go subscriberFuzz()
	}
}

func subscriber(spawned int) {
	ws := newWebSocket()

	for {
		path := fmt.Sprintf("/quick-fuzz/delayed/%d", rand.Intn(spawned))

		utilPause()

		err := ws.expect(
			fmt.Sprintf("/qio/on:1=\"%s\"", path),
			"/qio/callback/1:0={\"code\":200,\"data\":null}")
		if err != nil {
			log.Println("Subscribe failed:", err)
			continue
		}

		utilPause()
		err = ws.expect(
			fmt.Sprintf("/qio/off:2=\"%s\"", path),
			"/qio/callback/2:0={\"code\":200,\"data\":null}")
		if err != nil {
			log.Println("Unsubscribe failed:", err)
		}
	}
}

func subscriberFuzz() {
	ws := newWebSocket()

	for {
		utilPause()

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
