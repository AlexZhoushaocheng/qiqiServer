package router

import (
	"encoding/binary"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

const (
	//QUIT 停止
	QUIT = iota
)

var instance *Router

//Router Router
type Router struct {
	connQueue  chan *net.Conn
	signalChan chan int
}

//GetInstance 获取Router实例
func GetInstance() *Router {
	if nil == instance {
		instance = &Router{make(chan *net.Conn), make(chan int)}

	}
	return instance
}

//Push 添加一个新的连接到Router
func (router *Router) Push(conn *net.Conn) {
	router.connQueue <- conn
	//router.signalChan <- 1
}

//Run 启动
func (router *Router) Run() {
	go func() {
		for {
			select {
			case conn := <-router.connQueue:
				go router.handleConnection(*conn)
			case signal := <-router.signalChan:
				switch signal {
				case QUIT:
					return
				}
			}
		}
	}()
}

func packageData(data []byte) {
	fmt.Println(data)
}

//Stop 停止
func (router *Router) Stop() {
	router.signalChan <- QUIT
}

//TODO 最大长度限制 如果算出dataLength大于最大长度，则丢弃数据，并且回应消息错误
func (router *Router) handleConnection(conn net.Conn) {
	//defer conn.Close()
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
						//TODO 处理一个完整包
						packageData(data)
						dataBuf = dataBuf[dataLength+4:]
					}
				}
			} else { //正常情况,未发生粘包
				if length > 4 {
					dataLength := byteToInt32(readBuf[:4])
					if dataLength+4 == length {
						data = readBuf[4 : dataLength+4]
						//TODO 处理一个完整包
						packageData(data)
					} else if dataLength+4 < length {
						data = readBuf[4:dataLength]
						dataBuf = append(dataBuf, readBuf[dataLength+4:length]...)
						//TODO 处理一个完整包
						packageData(data)
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
