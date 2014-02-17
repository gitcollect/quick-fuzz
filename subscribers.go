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

		util_pause()
		qio.On(path, &fn)

		util_pause()
		qio.Off(path, nil)
	}
}

func subscriber_fuzz(spawned int) {
	qio := createClient()
	qio.Open()

	chCb := make(chan bool)
	for {
		util_pause()

		switch rand.Intn(4) {
		case 0:
			path := util_path_rand()
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
			path := util_path_valid_with_rand()
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
			path := util_path_rand()
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
			path := util_path_valid()
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
