package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
	"strconv"
)

func main() {
	var startValue = 0
	var lastRead = 0

	// setup listener
	pc_listener, err := net.ListenPacket("udp4", ":42069")
	if err != nil {
		panic(err)
	}
	buf_listener := make([]byte, 1024)

	isReceiving := true
	for isReceiving {
		time.Sleep(time.Second)
		i := listen(pc_listener, buf_listener)
		if i == 0 {
			fmt.Println("Sender is dead...\nQuit listening\nStart sending")
			startValue = lastRead;
			isReceiving = false
			pc_listener.Close()
		}else{
			lastRead = i
			fmt.Println(i)
		}
	}

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
	
	// spawn new listener
	err = exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go").Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	
	// primary work
	for {
		startValue++
		time.Sleep(time.Second)
		send(pc_broadcast, addr_broadcast, fmt.Sprint(startValue))
	}
}

func listen(pc net.PacketConn, buf []byte) (i int) {
	pc.SetReadDeadline(time.Now().Add(2 * time.Second))
	fmt.Println("Waiting for message...")
	n, _, err := pc.ReadFrom(buf)
	if err != nil {
		fmt.Println("Did not receive message")
		i = 0
	} else {
		j, err := strconv.Atoi(string(buf[:n]))
		if err != nil{
			fmt.Println("Could not convert str buffer to int")
		}
		i = j
	}
	return i
}

func send(pc net.PacketConn, addr *net.UDPAddr, message string) {
	_, err := pc.WriteTo([]byte(message), addr)
	if err != nil {
		panic(err)
	}
}
