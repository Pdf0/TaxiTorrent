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

	dataBase := make(map[string][]CentralProtocol.File)

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
		go processClient(connection, dataBase)
	}

}

func processClient(connection net.Conn, dataBase map[string][]CentralProtocol.File) {

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	s := new(CentralProtocol.SYN)
	util.DecodeToStruct(buffer[:mLen], s)

	fullAddr := net.JoinHostPort(s.Ip.String(), fmt.Sprintf("%d", s.Port))
	dataBase[fullAddr] = s.FileList

	fmt.Println("Received: ", *s)
	fmt.Println(dataBase)
	_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	connection.Close()

}
