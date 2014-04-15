package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	UuidChars = "abcdef123456790-"
)

var (
	FuzzMethods = [...]string{
		"GET",
		"OPTIONS",
		"PUT",
		"DELETE",
		"LAUGH",
		"PLAY",
		"HEAD",
		"DERP",
		"I don't understand HTTP",
	}
)

func httpReconnectors(spawn int) {
	for i := 0; i < spawn; i++ {
		go httpReconnector()
	}
}

func httpHeartbeaters(spawn int) {
	for i := 0; i < spawn; i++ {
		go httpHeartbeater()
	}
}

func httpFuzzes(spawn int) {
	for i := 0; i < spawn; i++ {
		go httpFuzz()
	}
}

func httpMultiRaces(spawn int) {
	for i := 0; i < spawn; i++ {
		go httpMultiRace()
	}
}

func httpRequest(uuid string, body string, connect bool) []byte {
	conn := ""
	if connect {
		conn = "&connect=true"
	}

	return []byte(fmt.Sprintf(
		"POST /?sid=%s%s HTTP/1.1\n"+
			"Content-Length: %d\n\n%s",
		uuid, conn, len(body), body))
}

func httpReconnector() {
	for {
		uuid := uuid.New()
		c := utilCreateSock()
		c.Write([]byte(httpRequest(uuid, "", true)))
		utilPause()
		c.Close()
	}
}

func httpHeartbeater() {
	buff := make([]byte, 1024)

	for {
		be_active := true
		uuid := uuid.New()

		c := utilCreateSock()
		c.Write(httpRequest(uuid, "", true))

		for be_active {
			if rand.Intn(2) == 1 {
				// Wait for the HTTP poll to finish
				l, _ := c.Read(buff)
				got := buff[:l]
				if !strings.Contains(string(got), "/qio/heartbeat:0=null") {
					log.Println("Error with resonse: %s", string(got))
					break
				}
			} else {
				// Force the first poll to finish
				c.Write(httpRequest(uuid, "/quick-fuzz/delayed:1234=null", false))

				// Response to first poll
				l, _ := c.Read(buff)
				got := buff[:l]
				if !strings.Contains(string(got), "/qio/heartbeat:0=null") {
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
		utilPause()
	}
}

func httpFuzzConnect() string {
	if rand.Intn(2) == 0 {
		return "&connect=true"
	} else {
		return ""
	}
}

func httpFuzzUuid() string {
	if rand.Intn(2) == 0 {
		uuid := make([]byte, rand.Intn(37))
		for i := range uuid {
			uuid[i] = UuidChars[rand.Intn(len(UuidChars))]
		}
		return string(uuid)
	} else {
		return uuid.New()
	}
}

func httpFuzz() {
	buff := make([]byte, 1024)

	for {
		bodyCnt := rand.Intn(10)
		httpMeth := FuzzMethods[rand.Intn(len(FuzzMethods))]
		connect := httpFuzzConnect()
		uuid := httpFuzzUuid()
		terminator := "\n\n"
		body := ""

		for i := 0; i < bodyCnt; i++ {
			body += utilRandomEvent()
		}

		bodyLen := len(body)
		if rand.Intn(2) == 0 && bodyLen > 0 {
			bodyLen = rand.Intn(bodyLen * 2)
		}

		if rand.Intn(8) == 0 {
			terminator = ""
		}

		msg := fmt.Sprintf(
			"%s /?sid=%s%s HTTP/1.1\n"+
				"Content-Length: %d%s"+
				"%s",
			httpMeth, uuid, connect, bodyLen, terminator, body)

		c := utilCreateSock()
		c.Write([]byte(msg))
		c.Read(buff)
		c.Close()
		utilPause()
	}
}

func httpMultiRaceRead(s net.Conn, recv chan string) {
	for {
		b := make([]byte, 2048)
		_, err := s.Read(b)

		if nErr, ok := err.(net.Error); ok && !nErr.Temporary() {
			return
		}

		recv <- string(b)
	}
}

func httpMultiRace() {
	for {
		got := 0
		cbId := 1
		uuid := uuid.New()
		socks := make([]net.Conn, rand.Intn(8)+1)
		pingsToSend := len(socks) * (rand.Intn(5) + 1)
		recv := make(chan string, pingsToSend)

		for i := range socks {
			s := utilCreateSock()
			s.Write([]byte(httpRequest(uuid, "", true)))
			socks[i] = s
		}

		for i := 0; i < pingsToSend; i++ {
			s := socks[i%len(socks)]
			ping := fmt.Sprintf("/qio/ping:%d=null", cbId)
			cbId++

			s.Write([]byte(httpRequest(uuid, ping, false)))
		}

		for _, s := range socks {
			go httpMultiRaceRead(s, recv)
		}

	recvFor:
		for got < pingsToSend {
			select {
			case b := <-recv:
				got += strings.Count(b, "/qio/callback/")

			case <-time.After(time.Second * 5):
				log.Println("Didn't get all ping callbacks after 5 seconds")
				break recvFor
			}
		}

		for _, s := range socks {
			s.Close()
		}

		utilPause()
	}
}
