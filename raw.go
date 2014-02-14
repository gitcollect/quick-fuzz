package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

func raws_fuzz_formatted(spawn int) {
	for i := 0; i < spawn; i++ {
		go raw_fuzz_formatted()
	}
}

func raws_fuzz_framed(spawn int) {
	for i := 0; i < spawn; i++ {
		go raw_fuzz_framed()
	}
}

func raws_fuzz_random(spawn int) {
	for i := 0; i < spawn; i++ {
		go raw_fuzz_random()
	}
}

func raw_fuzz_formatted() {
	buff := make([]byte, 8)
	bbuf := &bytes.Buffer{}

	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		c := createRawClient()

		bbuf.Reset()
		ev := fmt.Sprintf("%s:%s=%s", path(), path(), path())
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

func raw_fuzz_framed() {
	buff := make([]byte, 8)
	bbuf := &bytes.Buffer{}

	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		c := createRawClient()

		bbuf.Reset()
		path := path_rand()
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

func raw_fuzz_random() {
	buff := make([]byte, 8)

	for {
		<-time.After(time.Second * time.Duration(rand.Intn(SLEEP_MAX)))

		c := createRawClient()

		c.Write([]byte(path_rand()))

		_, err := c.Read(buff)
		if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
			log.Println("Temporary read error, server didn't close connection")
		}

		c.Close()
	}
}
