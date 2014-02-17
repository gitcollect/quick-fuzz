package main

import (
	"log"
	"net"
)

func insanes(spawn int) {
	for i := 0; i < spawn; i++ {
		go insane()
	}
}

func insane() {
	buff := make([]byte, 8)

	for {
		utilPause()
		c := utilCreateSock()

		c.Write([]byte(utilPath()))

		_, err := c.Read(buff)
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			log.Println("Temporary read error, server didn't close connection")
		}

		c.Close()
	}
}
