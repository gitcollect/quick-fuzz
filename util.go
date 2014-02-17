package main

import (
	"github.ihrint.com/quickio/quickigo"
	"log"
	"math/rand"
	"net"
	"net/url"
	"time"
)

const (
	SLEEP_MAX = 5
)

var (
	VALID_PATHS = []string{
		"/qio/ping",
		"/qio/on",
		"/qio/off",
		"/qio/callback",
		"/fuzzer/reject",
	}
)

const (
	ADDR = "unix:///tmp/quickio.sock"
)

func utilCreateClient() *quickigo.QuickIGo {
	return quickigo.New([]string{ADDR})
}

func utilCreateRawClient() net.Conn {
	for {
		c := utilCreateSock()

		_, err := c.Write([]byte(quickigo.HANDSHAKE))
		if err != nil {
			log.Println("createRawClient().Handshake_Write():", err)
			continue
		}

		res := make([]byte, len(quickigo.HANDSHAKE))
		_, err = c.Read(res)
		if err != nil {
			log.Println("createRawClient().Handshake_Recv():", err)
			continue
		}

		return c
	}
}

func utilCreateSock() net.Conn {
	for i := 0; ; i++ {
		url, err := url.Parse(ADDR)
		if err != nil {
			log.Fatal("createSock():", err)
			continue
		}

		c, err := net.Dial("unix", url.Path)
		if err != nil && i > 128 {
			log.Println("createSock().Dial():", err)
			continue
		}

		if c == nil {
			continue
		}

		return c
	}
}

func utilPause() {
	<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))
}

func utilPath() string {
	switch rand.Intn(3) {
	case 0:
		return utilPathRand()
	case 1:
		return utilPathValid()
	default:
		return utilPathValidWithRand()
	}
}

func utilPathRand() string {
	path := make([]byte, rand.Intn(128))
	for i := range path {
		path[i] = byte(rand.Intn(256))
	}

	return string(path)
}

func utilPathValid() string {
	return VALID_PATHS[rand.Intn(len(VALID_PATHS))]
}

func utilPathValidWithRand() string {
	path := utilPathValid()

	switch rand.Intn(3) {
	case 0:
		return utilPathRand() + path + utilPathRand()
	case 1:
		return utilPathRand() + path
	default:
		return path + utilPathRand()
	}
}
