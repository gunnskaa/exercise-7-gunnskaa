package main

import (
	"fmt"
	"time"
	"net"
)

func main() {
	var startCounter = 0

	// broadcast UDP setup
	pc_broadcast, err := net.ListenPacket("udp4", "") // automatic port number
	if err != nil {
		panic(err)
	}
	defer pc_broadcast.Close()

	addr_broadcast, err := net.ResolveUDPAddr("udp4", "255.255.255.255:42069")
	if err != nil {
		panic(err)
	}

	for {
		startCounter++
		time.Sleep(time.Second)
		send(pc_broadcast, addr_broadcast, fmt.Sprint(startCounter))
	}
}


func send(pc net.PacketConn, addr *net.UDPAddr, message string) {
	_, err := pc.WriteTo([]byte(message), addr)
	if err != nil {
		panic(err)
	}
}