package main

import (
	"TaxiTorrent/CentralProtocol"
	"TaxiTorrent/util"
	"fmt"
	"io"
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
		go processClient(connection, dataBase)
	}

}

func processClient(connection net.Conn, dataBase map[CentralProtocol.File][]string) {
	defer connection.Close()
	ip, port := util.GetTCPAddr(connection)
	fullAddr := net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port))
	fmt.Println("Node connected (" + fullAddr + ")")

	for {
		buffer := make([]byte, 1024)

		mLen, err := connection.Read(buffer)

		if err != nil {
			// Se o Node se disconectou, temos de o eliminar da base de dados
			if err == io.EOF {
				DisconnectNode(dataBase, fullAddr)
				fmt.Println("Node disconnected (" + fullAddr + ")")
			} else {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}
		if mLen == 0 {
			break
		}

		g := new(CentralProtocol.Central)
		util.DecodeToStruct(buffer[:mLen], g)

		switch g.PacketType {
		case "syn":
			s := new(CentralProtocol.SYN)
			if err := util.DecodeToStruct(g.Payload, s); err != nil {
				fmt.Println("Error decoding SYN packet:", err.Error())
				continue
			}
			InsertNodeInDataBase(dataBase, s.FileList, fullAddr)
			fmt.Println("Received: ", *s)

		case "update":
			u := new(CentralProtocol.Update)
			if err := util.DecodeToStruct(g.Payload, u); err != nil {
				fmt.Println("Error decoding Update packet:", err.Error())
				continue
			}
			UpdateNode(dataBase, u.FileList, fullAddr)
			fmt.Println("Updated Node (" + fullAddr + ")")

		case "get":

		}

		fmt.Println(dataBase)
	}
}

func UpdateNode(dataBase map[CentralProtocol.File][]string, files []CentralProtocol.File, fullAddr string) {
	DisconnectNode(dataBase, fullAddr)
	InsertNodeInDataBase(dataBase, files, fullAddr)
}

func DisconnectNode(dataBase map[CentralProtocol.File][]string, fullAddr string) {
	for file, list := range dataBase {
		dataBase[file] = util.RemoveStringFromList(list, fullAddr)
		if len(dataBase[file]) == 0 {
			delete(dataBase, file)
		}
	}
}

func InsertNodeInDataBase(dataBase map[CentralProtocol.File][]string, files []CentralProtocol.File, fullAddr string) {
	for _, file := range files {
		if _, ok := dataBase[file]; !ok {
			dataBase[file] = []string{fullAddr}
		} else {
			if !util.Contains(dataBase[file], fullAddr) {
				dataBase[file] = append(dataBase[file], fullAddr)
			}
		}
	}
}
