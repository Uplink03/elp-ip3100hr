package main

import (
	"fmt"
	"net"
	"os"
)

type PacketHandler struct {
	counter int
}

func (ph *PacketHandler) handle(packet *Packet) {
	ph.counter++
	filename := fmt.Sprintf("dumps/%d", ph.counter)
	ph.saveBin(filename+".bin", packet)
	ph.saveTxt(filename+".txt", packet)
	ph.savePayload(filename+".payload", packet)
}

func (ph *PacketHandler) saveBin(filename string, packet *Packet) {
	file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	data, _ := packet.MarshalBinary()
	file.Write(data)
}

func (ph *PacketHandler) saveTxt(filename string, packet *Packet) {
	file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	file.Write([]byte(packet.String()))
}

func (ph *PacketHandler) savePayload(filename string, packet *Packet) {
	file, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	file.Write(packet.Payload)
}

func main() {
	fmt.Println("Working directory:")
	fmt.Println(os.Getwd())

	if len(os.Args) < 3 {
		panic("Usage: protocol-dumper <CameraAddr:Port> <ListenAddr:Port>")
	}

	cameraAddr, err := net.ResolveTCPAddr("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}
	listenAddr, err := net.ResolveTCPAddr("tcp", os.Args[2])
	if err != nil {
		panic(err)
	}

	proxy := &Proxy{
		CameraAddr: cameraAddr,
		ListenAddr: listenAddr,
	}

	packetHandler := &PacketHandler{}
	proxy.Run(packetHandler.handle)
}
