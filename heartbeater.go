package main

import (
	"github.ihrint.com/quickio/quickigo"
)

func heartbeaters(spawn int) {
	for i := 0; i < spawn; i++ {
		go heartbeater()
	}
}

func heartbeater() {
	qio := quickigo.New([]string{ADDR})
	qio.Open()
	select {}
}
