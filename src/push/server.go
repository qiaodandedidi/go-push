package main

import (
	"librarys/log"
	"net"
	"sync"
)

type User struct {
	list map[int]net.Conn
	L    sync.RWMutex
}

var Users User = User{list: make(map[int]net.Conn, 10000)}

func Tcp(addr string) net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	log.FatalChk(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	log.FatalChk(err)
	return listener
}

func Read(listen *net.TCPListener) {
	tcpConn, err := listen.AcceptTCP()
	log.FatalChk(err)

}
