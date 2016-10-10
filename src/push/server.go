package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"librarys/log"
	"net"
	"sync"
)

type User struct {
	list map[int]*net.TCPConn
	L    sync.RWMutex
}

var UserList User = User{list: make(map[int]*net.TCPConn, 10000)}

func Tcp(addr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	log.FatalChk(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	log.FatalChk(err)
	return listener, nil
}

func Read(listen *net.TCPListener, ch chan string) {
	for {
		tcpConn, err := listen.AcceptTCP()
		log.FatalChk(err)
		ch <- "read"
		go read(tcpConn)
	}
}

func readHead(conn *net.TCPConn) (int8, int, string) {

	count := 4
	headBuf := make([]byte, 0)
	for {
		buffer := make([]byte, count)
		readLen, err := conn.Read(buffer)
		log.FatalChk(err)
		headBuf = append(headBuf, buffer[:readLen]...)
		count -= readLen
		if count == 0 {
			break
		}
	}
	var (
		version    int8
		contentLen int16
		cmd        string
	)

	fmt.Println(headBuf)
	err := binary.Read(bytes.NewReader(headBuf[:1]), binary.BigEndian, &version)
	log.FatalChk(err)
	err = binary.Read(bytes.NewReader(headBuf[2:]), binary.BigEndian, &contentLen)
	log.FatalChk(err)
	contentLenth := int(contentLen)
	buffers := bytes.NewBuffer(headBuf[1:2])
	cmd = buffers.String()
	return version, contentLenth, cmd
}
func read(conn *net.TCPConn) {
	err := conn.SetKeepAlive(true)
	log.FatalChk(err)
	version, contentLenth, cmd := readHead(conn)
	fmt.Println(version, contentLenth, cmd)
	contentBuf := make([]byte, 0)
	for {
		buffer := make([]byte, contentLenth)
		readLen, err := conn.Read(buffer)
		log.FatalChk(err)
		contentBuf = append(contentBuf, buffer[:readLen]...)
		contentLenth -= readLen
		fmt.Println(readLen, contentLenth)
		if contentLenth <= 0 {
			UserList.L.RLock()
			list := UserList.list
			UserList.L.RUnlock()
			for _, c := range list {
				go func(buf []byte) {
					c.Write(buf)
				}(contentBuf)
			}
			contentBuf = []byte{}
			_, contentLenth, _ = readHead(conn)
		}
	}
}

func Accept(listen *net.TCPListener, ch chan string) {
	i := 0
	for {
		i++
		tcpConn, err := listen.AcceptTCP()
		log.FatalChk(err)
		ch <- "write"
		UserList.L.Lock()
		UserList.list[i] = tcpConn
		UserList.L.Unlock()
	}
}

func Run(chOut, chIn chan string) {
	for {
		select {
		case v := <-chOut:
			fmt.Println(v)

		case v := <-chIn:
			fmt.Println(v)
		}
	}
}
