package main

import (
	config "GameServer/qiqiServer/GatewayServer/config"
	"GameServer/qiqiServer/GatewayServer/router"
	"net"

	log "github.com/sirupsen/logrus"
)

type server struct {
	port string
}

func (svr *server) run() {
	svr._init()
	router.GetInstance().Run()
	svr._startUp()
}

func (svr *server) _init() bool {
	ret := true
	svr.port = config.GetListenPort()
	return ret
}

func (svr *server) _startUp() {
	ln, err := net.Listen("tcp", ":"+svr.port)
	defer ln.Close()

	if err != nil {
		log.Error(err)
	} else {
		log.Info("listen success on " + ln.Addr().String())
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Error(err)
			} else {
				//go handleConnection(conn)
				router.GetInstance().Push(&conn)
			}
		}
	}
}

// func deal(data []byte) {
// 	data = append(data, 0)
// 	fmt.Println(string(data))
// }
