package main

import (
	"log"
	"math/rand"
	"net"
	"time"
)

func insanes(spawn int) {
	for i := 0; i < spawn; i++ {
		go insane()
	}
}

func insane() {
	buff := make([]byte, 8)

	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		c := createSock()

		c.Write([]byte(path()))

		_, err := c.Read(buff)
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			log.Println("Temporary read error, server didn't close connection")
		}

		c.Close()
	}
}
