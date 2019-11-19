package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/mm-uh/rpc_udp/src/util"
)

func main() {
	hostName := "localhost"
	portNum := "1053"

	service := hostName + ":" + portNum

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		log.Fatal(err)

	}

	log.Printf("Established connection to %s \n", service)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	defer conn.Close()

	// write a message to server
	rpcbase := &util.RPCBase{
		MethodName: "ExampleMethod",
		FirstArg:   "Joneeee",
	}
	toSend, err := json.Marshal(rpcbase)
	if err != nil {
		fmt.Println(err)
		return

	}

	message := []byte(string(toSend))

	_, err = conn.Write(message)

	if err != nil {
		log.Println("Errorrr: " + err.Error())

	}

	// receive message from server
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	var response util.ResponseRPC
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		fmt.Println("Error Unmarshaling response")
	}
	fmt.Println("UDP Server : ", addr)
	fmt.Println("Received from UDP server : ", response.Response)

}
