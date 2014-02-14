package main

import (
	"fmt"
	"github.ihrint.com/quickio/quickigo"
	"log"
	"math/rand"
	"time"
)

func subscribers(spawn int) {
	for i := 0; i < spawn; i++ {
		go subscriber(spawn)
	}
}

func subscribers_fuzz(spawn int) {
	for i := 0; i < spawn; i++ {
		go subscriber_fuzz(spawn)
	}
}

func subscriber(spawned int) {
	qio := createClient()
	qio.Open()

	var fn quickigo.SubCb = func(data interface{}) {}

	for {
		path := fmt.Sprintf("/fuzzer/delayed/%d", rand.Intn(spawned))

		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))
		qio.On(path, &fn)

		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))
		qio.Off(path, nil)
	}
}

func subscriber_fuzz(spawned int) {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		switch rand.Intn(4) {
		case 0:
			path := path_rand()
			qio.Send("/qio/on", path,
				func(_ interface{}, _ quickigo.ServerCbFn, code int, _ string) {
					if code == quickigo.CODE_OK {
						log.Println("Subscribe fuzzer: invalid event not "+
							"rejected:", path)
						qio.Send("/qio/off", path, nil)
					}
					chCb <- true
				})

		case 1:
			path := path_valid_with_rand()
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
			path := path_rand()
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
			path := path_valid()
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
