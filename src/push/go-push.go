package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
)

/*
#include <unistd.h>
*/
import "C"

type user struct {
	list map[int]net.Conn
	L    sync.RWMutex
}

var userList user = user{list: make(map[int]net.Conn)}

func main() {
	C.daemon(1, 0)
	runtime.GOMAXPROCS(runtime.NumCPU())
	addrOut, err1 := net.ResolveTCPAddr("tcp", ":8011")
	checkError(err1)
	addrIn, err2 := net.ResolveTCPAddr("tcp", ":8088")
	checkError(err2)
	listenerOut, errOut := net.ListenTCP("tcp", addrOut)
	checkError(errOut)
	listenerIn, errIn := net.ListenTCP("tcp", addrIn)
	checkError(errIn)
	fmt.Println("start success")
	chOut := make(chan string)
	chIn := make(chan string)

	go func() {
		i := 0
		for {
			i++
			conn, err := listenerOut.AcceptTCP()
			checkError(err)
			userList.L.Lock()
			userList.list[i] = conn
			userList.L.Unlock()
			go handleRead(conn, chOut)
		}
	}()
	go func() {
		for {
			conn, err := listenerIn.AcceptTCP()
			checkError(err)
			go handleWrite(conn, chIn)
		}
	}()
	for {
		select {
		case v := <-chOut:
			fmt.Println(v)
		case v := <-chIn:
			fmt.Println(v)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func handleRead(conn *net.TCPConn, ch chan string) {
	ch <- "read"
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err == nil {
			fmt.Println(buf[:n])
		}
	}
}
func handleWrite(conn *net.TCPConn, ch chan string) {
	ch <- "write"
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err == nil {
			userList.L.RLock()
			list := userList.list
			userList.L.RUnlock()
			for _, c := range list {
				c.Write(buf[:n])
			}
		}
	}
}
