package main

import (
	"fmt"
	"time"
	"os/exec"
	"net"
)


func main() {
	var global_counter = 0

	// listener setup
	listenChan := make(chan bool)
	pc_listener, err := net.ListenPacket("udp4", ":42069")
	if err != nil {
		panic(err)
	}
	defer pc_listener.Close()
	buf_listener := make([]byte, 1024)

	// launch a new instance
	err = exec.Command("cmd", "/C", "start", "powershell", "go", "run", "pheonix.go").Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	
	// go watchdog to kill listener
	go listen(pc_listener, buf_listener, listenChan)
	var secondsWaited = 0
	for{
		select{
		case <-listenChan:
			// save variable
			global_counter++
			fmt.Println(global_counter)
			time.Sleep(time.Second)
		default:
			if secondsWaited < 3{
				time.Sleep(time.Second)
				secondsWaited++
			}else{
				secondsWaited = 0
				// launch a new instance
				err = exec.Command("cmd", "/C", "start", "powershell", "go", "run", "pheonix.go").Run()
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

func listen(pc net.PacketConn, buf []byte, listenChan chan<- bool){
	for {
		_, _, err := pc.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		listenChan <- true		
		//fmt.Printf("%s sent this: %s\n", addr, buf[:n])
	}
}