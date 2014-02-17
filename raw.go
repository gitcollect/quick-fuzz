package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func rawsFuzzFormatted(spawn int) {
	for i := 0; i < spawn; i++ {
		go rawFuzzFormatted()
	}
}

func rawsFuzzFramed(spawn int) {
	for i := 0; i < spawn; i++ {
		go rawFuzzFramed()
	}
}

func rawsFuzzRandom(spawn int) {
	for i := 0; i < spawn; i++ {
		go rawFuzzRandom()
	}
}

func rawFuzzFormatted() {
	buff := make([]byte, 8)
	bbuf := &bytes.Buffer{}

	for {
		utilPause()
		c := utilCreateRawClient()

		bbuf.Reset()
		ev := fmt.Sprintf("%s:%s=%s", utilPath(), utilPath(), utilPath())
		binary.Write(bbuf, binary.BigEndian, len(ev))
		bbuf.WriteString(ev)

		bbuf.WriteTo(c)

		_, err := c.Read(buff)
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			log.Println("Temporary read error, server didn't close connection")
		}

		c.Close()
	}
}

func rawFuzzFramed() {
	buff := make([]byte, 8)
	bbuf := &bytes.Buffer{}

	for {
		utilPause()
		c := utilCreateRawClient()

		bbuf.Reset()
		path := utilPathRand()
		binary.Write(bbuf, binary.BigEndian, len(path))
		bbuf.WriteString(path)

		bbuf.WriteTo(c)

		_, err := c.Read(buff)
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			log.Println("Temporary read error, server didn't close connection")
		}

		c.Close()
	}
}

func rawFuzzRandom() {
	buff := make([]byte, 8)

	for {
		utilPause()
		c := utilCreateRawClient()

		c.Write([]byte(utilPathRand()))

		_, err := c.Read(buff)
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			log.Println("Temporary read error, server didn't close connection")
		}

		c.Close()
	}
}
