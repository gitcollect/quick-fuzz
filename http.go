package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

const (
	UuidFuzzChars = "ABCDEFabcdef123456790ijklmnaoeqwe*&^(#"
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
		uuid := uuid.NewV4().String()
		c := utilCreateSock()
		c.Write([]byte(httpRequest(uuid, "", true)))
		utilPause()
		c.Close()
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
		uuid := make([]byte, rand.Intn(33))
		for i := range uuid {
			uuid[i] = UuidFuzzChars[rand.Intn(len(UuidFuzzChars))]
		}
		return string(uuid)
	} else {
		return uuid.NewV4().String()
	}
}

func httpFuzz() {
	buff := make([]byte, 1024)

	for {
		utilPause()

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
		utilPause()

		got := 0
		cbId := 1
		uuid := uuid.NewV4().String()
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
	}
}
