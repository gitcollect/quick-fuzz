package main

import (
	"code.google.com/p/go-uuid/uuid"
	"math/rand"
	"net"
)

func insanes(spawn int) {
	for i := 0; i < spawn; i++ {
		go insane()
	}
}

func insane() {
	buff := make([]byte, 1024)

	for {
		utilPause()

		var c net.Conn
		uuid := uuid.New()

		switch rand.Intn(4) {
		case 0:
			c = utilWebSocketClient(buff)

		case 1:
			c = utilCreateSock()
			c.Write([]byte(httpRequest(uuid, "", true)))

		case 2:
			c = utilCreateRawClient()

		case 3:
			c = utilCreateSock()
			c.Write([]byte("<policy-file-request/>"))
		}

		switch rand.Intn(2) {
		case 0:
			c.Write([]byte(httpRequest(uuid, utilRandomEvent(), false)))

		case 1:
			c.Write([]byte(utilRandomEvent()))
		}

		c.Close()
	}
}
