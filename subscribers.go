package main

import (
	"fmt"
	"github.ihrint.com/quickio/quickigo"
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
		go subscriberFuzz(spawn)
	}
}

func subscriber(spawned int) {
	qio := utilCreateClient()
	qio.Open()

	var fn quickigo.SubCb = func(data interface{}) {}

	for {
		path := fmt.Sprintf("/fuzzer/delayed/%d", rand.Intn(spawned))

		utilPause()
		qio.On(path, &fn)

		utilPause()
		qio.Off(path, nil)
	}
}

func subscriberFuzz(spawned int) {
	qio := utilCreateClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		utilPause()

		switch rand.Intn(4) {
		case 0:
			path := utilPathRand()
			qio.Send("/qio/ron", path,
				func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
					if code == quickigo.CODE_OK {
						log.Println("Subscribe fuzzer: invalid event not "+
							"rejected:", path)
						qio.Send("/qio/off", path, nil)
					}
					chCb <- true
				})

		case 1:
			path := utilPathValidWithRand()
			qio.Send("/qio/on", path,
				func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
					if code == quickigo.CODE_OK {
						log.Println("Subscribe fuzzer: invalid event not "+
							"rejected:", path)
						qio.Send("/qio/off", path, nil)
					}
					chCb <- true
				})

		case 2:
			path := utilPathRand()
			qio.Send("/qio/off", path,
				func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
					if code != quickigo.CODE_NOT_FOUND {
						log.Printf("Event not rejected with "+
							"CODE_NOT_FOUND (got %d): %s\n",
							code, path)
					}
					chCb <- true
				})

		case 3:
			path := utilPathValid()
			qio.Send("/qio/off", path,
				func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
					if code != quickigo.CODE_OK {
						log.Printf("Event not accepted with CODE_OK (got %d): %s\n",
							code, path)
					}
					chCb <- true
				})
		}

		<-chCb
	}
}
