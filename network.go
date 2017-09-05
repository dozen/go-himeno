package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
	DefaultProtocol      = "tcp"
	DefaultListenAddr    = "0.0.0.0:22122"
	DefaultRPCListenAddr = "0.0.0.0:22123"

	MethodCLIENT = "client"
	MethodSERVER = "server"

	MessageTrue  = "T"
	MessageFalse = "F"
)

func init() {
	//Control Server
	ctrlSrvCh := make(chan error)
	ServeTCP(DefaultProtocol, DefaultRPCListenAddr, controlServerHandler, ctrlSrvCh)

	select {
	case err := <-ctrlSrvCh:
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}
}

type connHandler func(conn *net.TCPConn)

func controlServerHandler(conn *net.TCPConn) {
	b := make([]byte, 20)
	size, err := conn.Read(b)
	if err != nil {
		fmt.Errorf(err.Error())
		conn.Write([]byte(MessageFalse))
		conn.Write([]byte(err.Error()))
	}

	receiveMsg := string(size[:size])
	switch receiveMsg {
	case "stats":
		// なんか状態とかを取得できる
	case "member":
		// メンバー一覧を取得する
	case "start":
		// スタートの合図
	}
}

func ServeTCP(protocol, listenAddr string, handler connHandler, errCh chan error) {
	lAddr, err := net.ResolveTCPAddr(protocol, listenAddr)
	if err != nil {
		errCh <- err
		return
	}
	listener, err := net.ListenTCP(protocol, lAddr)
	if err != nil {
		errCh <- err
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Errorf(err.Error())
			continue
		}
		go handler(conn)
	}
}

func receive(conn io.Reader) {
	msg := make([]byte, 6)
	conn.Read(msg)

	if bytes.Compare(msg, methodClient()) {

	}

}

func methodClient() []byte {
	return []byte(MethodCLIENT)
}

func methodServer() []byte {
	return []byte(MethodSERVER)
}
