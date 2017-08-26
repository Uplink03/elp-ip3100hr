package main

import (
	"fmt"
	"io"
	"net"
)

type Proxy struct {
	CameraAddr *net.TCPAddr
	ListenAddr *net.TCPAddr
}

type ProxyHandler func(packet *Packet)

func (proxy *Proxy) Run(handler ProxyHandler) {
	listener, err := net.ListenTCP("tcp", proxy.ListenAddr)
	if err != nil {
		panic(fmt.Sprintf("Could not listen: %+v\n", err))
	}

	for {
		appConn, err := listener.AcceptTCP()
		if err != nil {
			panic(fmt.Sprintf("Could not accept: %+v", err))
		}
		fmt.Printf("Accepted: %+v\n", appConn)

		go proxy.cameraProxyConnection(appConn, handler)
	}
}

func (proxy *Proxy) cameraProxyConnection(appConn *net.TCPConn, handler ProxyHandler) {
	cameraConn, err := net.DialTCP("tcp", nil, proxy.CameraAddr)
	if err != nil {
		panic(fmt.Sprintf("Couldn't dial: %+v", err))
	}
	fmt.Printf("Connected: %+v\n", cameraConn)

	cameraRecv := make(chan *Packet)
	appRecv := make(chan *Packet)

	go packetReader(cameraRecv, cameraConn)
	go packetReader(appRecv, appConn)

forloop:
	for {
		select {
		case appData, ok := <-appRecv:
			if !ok {
				break forloop
			}
			handler(appData)
			headerData, _ := appData.Header.MarshalBinary()
			cameraConn.Write(headerData)
			cameraConn.Write(appData.Payload)
		case cameraData, ok := <-cameraRecv:
			if !ok {
				break forloop
			}
			handler(cameraData)
			headerData, _ := cameraData.Header.MarshalBinary()
			appConn.Write(headerData)
			appConn.Write(cameraData.Payload)
		}
	}

	fmt.Println("Disconnected")
}

func packetReader(recv chan<- *Packet, conn *net.TCPConn) {
	defer func() {
		close(recv)
		conn.Close()
	}()

	for {
		headerData := make([]byte, HeaderLength)
		_, err := io.ReadFull(conn, headerData)
		if err != nil {
			return
		}

		packet := &Packet{}
		err = packet.Header.UnmarshalBinary(headerData)
		if err != nil {
			return
		}

		packet.Payload = make([]byte, packet.PayloadLength)
		_, err = io.ReadFull(conn, packet.Payload)
		if err != nil {
			return
		}

		recv <- packet
	}
}
