package main

import (
	"encoding/json"
	"fmt"
	serverRpc "github.com/mm-uh/rpc_udp/src/server"
	"github.com/mm-uh/rpc_udp/src/util"
	"log"
	"net"
)

var exit1 = make(chan bool)
var exit2 = make(chan bool)

type Handler int

func (h *Handler) Ping(i string, j float64) string {
	fmt.Println("i -> ", i)
	fmt.Println("j -> ", j)

	return "Pong"
}

func (h *Handler) WithTwo(i, j string) string {
	fmt.Println("i -> ", i)
	fmt.Println("j -> ", j)
	return "Mauricio es llegua"
}

func main() {
	var h Handler
	server := serverRpc.NewServer(h, ":1053")
	// listen to incoming udp packets
	var exited = make(chan bool)
	go server.ListenServer(exited)
	go client(1)
	//go client(2)
	if s := <-exited; s {
		// Handle Error in method
		fmt.Println("We get an error listen server")
		return
	}
	<-exit1
	<-exit2
}

func client(method int16) {
	hostName := "localhost"
	portNum := "1053"

	service := hostName + ":" + portNum

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	if err != nil {
		log.Fatal(err)

	}

	log.Printf("Established connection to %s \n", service)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	defer conn.Close()

	rpcbase := &util.RPCBase{
		MethodName: "",
	}
	some := make([]interface{}, 0)
	switch method {
	case 1:
		{
			rpcbase.MethodName = "Ping"
			some = append(some, int(1))
			some = append(some, string("45"))
		}
	case 2:
		{
			rpcbase.MethodName = "WithTwo"
			some = append(some, "Ping")
			some = append(some, "Ping")
		}
	}
	rpcbase.Args = some

	toSend, err := json.Marshal(rpcbase)
	if err != nil {
		fmt.Println(err)
		return

	}

	message := []byte(string(toSend))

	for i := 0; ; i++ {
		_, err = conn.Write(message)

		if err != nil {
			log.Println("Errorrr: " + err.Error())
			break
		}

		// receive message from server
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)

		var response util.ResponseRPC
		err = json.Unmarshal(buffer[:n], &response)
		if err != nil {
			fmt.Println("Error Unmarshaling response")
			break
		}
		fmt.Println("ITERATION ", i)
		fmt.Println("UDP Server : ", addr)
		fmt.Println("Received from UDP server : ", response.Response)

	}

	switch method {
	case 1:
		{
			exit1 <- true
		}
	case 2:
		{
			exit2 <- true
		}

	}
}
