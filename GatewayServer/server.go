package main

import (
	config "GameServer/GatewayServer/config"
	"encoding/binary"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

type server struct {
	port string
}

func (svr *server) run() {
	svr._init()
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
				go handleConnection(conn)
			}
		}
	}
}

func deal(data []byte) {
	data = append(data, 0)
	fmt.Println(string(data))
}

//TODO 最大长度限制 如果算出dataLength大于最大长度，则丢弃数据，并且回应消息错误
func handleConnection(conn net.Conn) {
	defer conn.Close()
	readBuf := make([]byte, 1024) //socket读缓冲
	data := make([]byte, 0)       //完整的数据包
	dataBuf := make([]byte, 0)    //有粘包时缓冲
	for {
		length, err := conn.Read(readBuf)
		if nil == err {
			if len(dataBuf) > 0 { //粘包处理
				dataBuf = append(dataBuf, readBuf[:length]...)

				if len(dataBuf) > 4 {
					dataLength := byteToInt32(dataBuf[:4])
					if len(dataBuf) >= dataLength+4 {
						data = dataBuf[4 : dataLength+4]
						deal(data)
						dataBuf = dataBuf[dataLength+4:]
					}
				}
			} else { //正常情况,未发生粘包
				if length > 4 {
					dataLength := byteToInt32(readBuf[:4])
					if dataLength+4 == length {
						data = readBuf[4 : dataLength+4]
						//TODO 处理一个完整包
						deal(data)
					} else if dataLength+4 < length {
						data = readBuf[4:dataLength]
						dataBuf = append(dataBuf, readBuf[dataLength+4:length]...)
						//TODO 处理一个完整包
						deal(data)
					} else {
						dataBuf = append(dataBuf, readBuf[:length]...)
					}
				} else {
					dataBuf = append(dataBuf, readBuf[:length]...)
				}
			}
		} else {
			log.Error(err)
		}
	}
}

func byteToInt32(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}
