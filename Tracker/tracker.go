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

	dataBase := make(map[CentralProtocol.File][]string)

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

func processClient(connection net.Conn, dataBase map[CentralProtocol.File][]string) {
	defer connection.Close()

	for {
		buffer := make([]byte, 1024)
		mLen, err := connection.Read(buffer)

		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}

		if mLen == 0 {
			break
		}

		s := new(CentralProtocol.SYN)
		util.DecodeToStruct(buffer[:mLen], s)

		fullAddr := net.JoinHostPort(s.Ip.String(), fmt.Sprintf("%d", s.Port))
		for _, file := range s.FileList {
			if _, ok := dataBase[file]; !ok {
				dataBase[file] = []string{fullAddr}
			} else {
				if !util.Contains(dataBase[file], fullAddr) {
					dataBase[file] = append(dataBase[file], fullAddr)
				}
			}
		}
		fmt.Println("Received: ", *s)
		fmt.Println(dataBase)
	}
}
