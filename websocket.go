package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	WebSocketHandshake     = "\x81\x89\x00\x00\x00\x00/qio/ohai"
	WebSocketHandshakeResp = "\x81\x09/qio/ohai"
	WebSocketHeaders       = "GET / HTTP/1.1\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: JF+JVs2N4NAX39FAAkkdIA==\r\n" +
		"Sec-WebSocket-Protocol: quickio\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"
)

var (
	WebSocketMask   = []byte{0, 0, 0, 0}
	ErrFrameTooLong = errors.New("WebSocket client does not support frames >4096 bytes")
	ErrNotConnected = errors.New("WebSocket is not connected, so recv will never get anything")
)

type webSocket struct {
	buff []byte
	conn net.Conn
}

func newWebSocket() *webSocket {
	return &webSocket{
		buff: make([]byte, 4096),
	}
}

func (ws *webSocket) open() {
	ws.close()
	ws.conn = utilWebSocketClient(ws.buff)
}

func (ws *webSocket) close() {
	if ws.conn != nil {
		ws.conn.Close()
		ws.conn = nil
	}
}

func (ws *webSocket) send(msg string) (err error) {
	if ws.conn == nil {
		ws.open()
	}

	header_len := 6 // always at least 6: 2 opening btyes + mask
	size := len(msg)

	ws.buff[0] = 0x81

	if size <= 125 {
		ws.buff[1] = 0x80 | byte(size)
	} else if size <= 65536 {
		ws.buff[1] = 0x80 | 126
		binary.BigEndian.PutUint16(ws.buff[2:], uint16(size))
		header_len += 2
	} else {
		err = ErrFrameTooLong
		return
	}

	if (header_len + size) > len(ws.buff) {
		err = ErrFrameTooLong
		return
	}

	copy(ws.buff[header_len-4:], WebSocketMask)

	msg_start := ws.buff[header_len:]
	for i := 0; i < size; i++ {
		msg_start[i] = WebSocketMask[i&3] ^ []byte(msg)[i]
	}

	_, err = ws.conn.Write(ws.buff[:header_len+size])
	if err != nil {
		ws.close()
	}
	return
}

func (ws *webSocket) recv(msg string, prefix bool) (err error) {
	if ws.conn == nil {
		return ErrNotConnected
	}

	_, err = ws.conn.Read(ws.buff)
	if err != nil {
		ws.close()
		return
	}

	opcode := ws.buff[0] & 0x0f
	if opcode != 1 {
		err = errors.New(fmt.Sprintf("Expected opcode of 1, got opcode %d",
			opcode))
		return
	}

	header_len := 2

	size := int(ws.buff[1] & 0x7f)
	if size == 126 {
		size = int(binary.BigEndian.Uint16(ws.buff[2:]))
		header_len += 2
	} else if size == 127 {
		err = ErrFrameTooLong
		return
	}

	resp := string(ws.buff[header_len : header_len+size])

	if prefix {
		if !strings.HasPrefix(resp, msg) {
			err = errors.New(fmt.Sprintf("Expected \"%s\" to have prefix \"%s\"", resp, msg))
		}
	} else {
		if resp != msg {
			err = errors.New(fmt.Sprintf("Expected \"%s\", got \"%s\"", msg, resp))
		}
	}

	return
}

func (ws *webSocket) expect(msg string, resp string) (err error) {
	err = ws.send(msg)
	if err != nil {
		return
	}

	err = ws.recv(resp, false)
	return
}

func (ws *webSocket) expectPrefix(msg string, resp string) (err error) {
	err = ws.send(msg)
	if err != nil {
		return
	}

	err = ws.recv(resp, true)
	return
}
