package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/iheartradio/quickigo"
)

const (
	// Addr         = "ws://localhost:8080"
	Addr         = "unix:///tmp/quickio.sock"
	QioHandshake = "/qio/ohai"
	SleepMax     = 5
	UuidChars    = "ABCDEFabcdef123456790"
)

var (
	ValidPaths = []string{
		"/qio/ping",
		"/qio/on",
		"/qio/off",
		"/qio/callback",
		"/qio/hostname",
		"/quick-fuzz/reject",
	}
)

func utilCreateClient() *quickigo.QuickIGo {
	return quickigo.New(Addr)
}

func utilCreateRawClient() net.Conn {
	res := make([]byte, len(QioHandshake))

	for {
		c := utilCreateSock()

		_, err := c.Write([]byte(QioHandshake))
		if err != nil {
			log.Println("createRawClient().Handshake_Write():", err)
			continue
		}

		_, err = c.Read(res)
		if err != nil {
			log.Println("createRawClient().Handshake_Recv():", err)
			continue
		}

		return c
	}
}

func utilWebSocketClient(buff []byte) net.Conn {
	for {
		c := utilCreateSock()

		_, err := c.Write([]byte(WebSocketHeaders))
		if err != nil {
			log.Println("Failed to send websocket handshake:", err)
			continue
		}

		got, err := c.Read(buff)
		if err != nil {
			log.Println("Failed to read websocket handshake:", err)
			continue
		}

		resp := string(buff[:got])
		if !strings.Contains(resp, "Nf+/kB4wxkn+6EPeanngB3VZNwU=") {
			log.Println("Invalid websocket handshake response:", resp)
			continue
		}

		_, err = c.Write([]byte(WebSocketHandshake))
		if err != nil {
			log.Println("Failed to write QIO handshake:", err)
			continue
		}

		_, err = c.Read(buff)
		if err != nil {
			log.Println("Failed to read QIO handshake response:", err)
			continue
		}

		return c
	}
}

func utilCreateSock() net.Conn {
	for i := 0; ; i++ {
		var c net.Conn

		url, err := url.Parse(Addr)
		if err != nil {
			log.Fatal("createSock():", err)
			continue
		}

		host := url.Host
		if url.Scheme == "unix" {
			c, err = net.Dial("unix", url.Path)
		} else if url.Scheme == "wss" {
			if !strings.Contains(host, ":") {
				host += ":443"
			}

			c, err = tls.Dial("tcp", host, nil)
		} else {
			if !strings.Contains(host, ":") {
				host += ":80"
			}

			c, err = net.Dial("tcp", host)
		}

		if err != nil {
			if i > 128 {
				log.Println("createSock().Dial():", err)
			}
			continue
		}

		return c
	}
}

func utilPause() {
	<-time.After(time.Second * time.Duration(rand.Intn(SleepMax)))
}

func randRune() rune {
	for {
		i := rand.Intn(unicode.MaxRune)
		r := rune(i)
		if i > 0 && r != ':' && r != '=' {
			return r
		}
	}
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
	path := make([]rune, rand.Intn(128))
	for i := range path {
		path[i] = randRune()
	}

	return string(path)
}

func utilPathValid() string {
	return ValidPaths[rand.Intn(len(ValidPaths))]
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

func utilRandomEvent() string {
	return fmt.Sprintf("%s:%d=%s\n",
		utilPath(),
		rand.Intn(100000),
		utilPathRand())
}
