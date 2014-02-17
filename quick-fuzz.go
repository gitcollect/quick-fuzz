package main

import (
	"github.ihrint.com/quickio/quickigo"
	"log"
	"net"
	"net/url"
	"runtime"
)

const (
	ADDR = "unix:///tmp/quickio.sock"
)

func createClient() *quickigo.QuickIGo {
	return quickigo.New([]string{ADDR})
}

func createRawClient() net.Conn {
	for {
		c := createSock()

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

func createSock() net.Conn {
	for i := 0; ; i++ {
		url, err := url.Parse(ADDR)
		if err != nil {
			log.Fatal("createSock():", err)
			continue
		}

		c, err := net.Dial("unix", url.Path)
		if err != nil && i > 10 {
			log.Println("createSock().Dial():", err)
			continue
		}

		if c == nil {
			continue
		}

		return c
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	callbacks(15000)
	callbacks_fuzz_chain(15000)
	callbacks_fuzz_recv(15000)
	callbacks_fuzz_send(15000)
	heartbeaters(20000)
	insanes(50000)
	raws_fuzz_formatted(15000)
	raws_fuzz_framed(15000)
	raws_fuzz_random(15000)
	reconnectors(2000)
	subscribers(15000)
	subscribers_fuzz(15000)

	select {}
}
