package main

import (
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

func wsHeartbeaters(spawn int) {
	for i := 0; i < spawn; i++ {
		go wsHeartbeater()
	}
}

func httpHeartbeaters(spawn int) {
	for i := 0; i < spawn; i++ {
		go httpHeartbeater()
	}
}

func wsHeartbeater() {
	qio := utilCreateClient()
	qio.Open()
	select {}
}

func httpHeartbeater() {
	buff := make([]byte, 1024)

	for {
		be_active := true
		uuid := utilUuid()

		c := utilCreateSock()
		c.Write(httpRequest(uuid, "", true))

		for be_active {
			utilPause()

			if rand.Intn(2) == 1 {
				// Wait for the HTTP poll to finish
				l, _ := c.Read(buff)
				got := buff[:l]
				if !strings.Contains(string(got), "Content-Length: 0") {
					log.Println("Error with resonse:", string(got))
					break
				}
			} else {
				// Force the first poll to finish
				c.Write(httpRequest(uuid, "/quick-fuzz/delayed:1234=null", false))

				// Response to first poll
				l, _ := c.Read(buff)
				got := buff[:l]
				if !strings.Contains(string(got), "Content-Length: 0") {
					log.Println("Error with resonse:", string(got))
					break
				}

				if !strings.Contains(string(got), "/qio/callback/1234") {
					// Callback to /quick-fuzz/delayed, fired from second poll
					l, _ = c.Read(buff)
					got = buff[:l]
					if !strings.Contains(string(got), "/qio/callback/1234") {
						log.Println("Didn't get callback from /quick-fuzz/delayed:", string(got))
						break
					}
				}
			}

			// At this point, there is no active poll on the server

			switch rand.Intn(3) {
			case 0:
				// Wait for the surrogate to timeout and verify 403
				be_active = false

				<-time.After(time.Second * 21)

				c.Write(httpRequest(uuid, "", false))
				l, _ := c.Read(buff)
				got := buff[:l]
				if !strings.Contains(string(got), "403 Forbidden") {
					log.Println("Did not get 403 after not polling on surrogate:", string(got))
				}

			case 1:
				// Wait for the socket to close
				be_active = false

				c.SetReadDeadline(time.Now().Add(65 * time.Second))
				_, err := c.Read(buff)

				if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
					log.Println("Server did not close HTTP connection after 65 seconds")
				}

			default:
				// Just poll again
				c.Write(httpRequest(uuid, "", false))
			}
		}

		c.Close()
	}
}
