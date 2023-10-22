package main

import (
	"TaxiTorrent/CentralProtocol"
	"TaxiTorrent/util"
	"fmt"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "10000"
	SERVER_TYPE = "tcp"
)

func main() {

	fmt.Println("Server Running...")
	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
	}

	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}
		fmt.Println("client connected")
		go processClient(connection)
	}

}

func processClient(connection net.Conn) {

	s := new(CentralProtocol.SYN)
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	util.DecodeToStruct(buffer[:mLen], s)
	fmt.Println("Received: ", *s)
	_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	connection.Close()

}
