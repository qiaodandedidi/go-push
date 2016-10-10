package main

import (
	"conf"
	"fmt"
	"librarys/log"
	"net"
	"runtime"
	"sync"
)

type user struct {
	list map[int]net.Conn
	L    sync.RWMutex
}

var userList user = user{list: make(map[int]net.Conn)}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var (
		outAddr, inAddr string
	)
	if _, ok := conf.Conf["outterAddr"].(string); !ok {
		return
	}
	outAddr = conf.Conf["outterAddr"].(string)
	listenerOut, err := Tcp(outAddr)
	log.FatalChk(err)
	if _, ok := conf.Conf["innerAddr"].(string); !ok {
		return
	}
	inAddr = conf.Conf["innerAddr"].(string)
	listenerIn, err := Tcp(inAddr)
	log.FatalChk(err)
	chOut := make(chan string)
	chIn := make(chan string)
	go Read(listenerIn, chIn)
	go Accept(listenerOut, chOut)
	fmt.Println("start success")
	Run(chOut, chIn)
}
