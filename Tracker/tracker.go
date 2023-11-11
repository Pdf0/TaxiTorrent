package main

import (
	"TaxiTorrent/CentralProtocol"
	"TaxiTorrent/util"
	"fmt"
	"io"
	"net"
	"reflect"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "10000"
	SERVER_TYPE = "tcp"
	BLOCKSIZE   = 256
)

type Seeder struct {
	Ip     net.IP
	Port   uint
	Blocks []string
}

func main() {

	dataBase := make(map[string][]Seeder)

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

func processClient(connection net.Conn, dataBase map[string][]Seeder) {
	defer connection.Close()
	ip, port := util.GetTCPRemoteAddr(connection)
	fullAddr := net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port))
	fmt.Println("Node connected (" + fullAddr + ")")

	for {
		buffer := make([]byte, 1024)

		mLen, err := connection.Read(buffer)

		if err != nil {
			if err == io.EOF {
				DisconnectNode(dataBase, ip, port)
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
			InsertNodeInDataBase(dataBase, s.FileList, s.Ip, s.Port)
			fmt.Println("Received: ", *s)

		case "update":
			u := new(CentralProtocol.Update)
			if err := util.DecodeToStruct(g.Payload, u); err != nil {
				fmt.Println("Error decoding Update packet:", err.Error())
				continue
			}
			UpdateNode(dataBase, u.FileList, ip, port)
			fmt.Println("Updated Node (" + fullAddr + ")")

		case "get":

		}

		fmt.Println(dataBase)
	}
}

func UpdateNode(dataBase map[string][]Seeder, files []CentralProtocol.File, ip net.IP, port uint) {
	DisconnectNode(dataBase, ip, port)
	InsertNodeInDataBase(dataBase, files, ip, port)
}

func DisconnectNode(dataBase map[string][]Seeder, ip net.IP, port uint) {
	for file, list := range dataBase {
		dataBase[file] = RemoveNodeFromList(list, ip)
		if len(dataBase[file]) == 0 {
			delete(dataBase, file)
		}
	}
}

func InsertNodeInDataBase(dataBase map[string][]Seeder, files []CentralProtocol.File, ip net.IP, port uint) {
	for _, file := range files {
		node := Seeder{
			ip,
			port,
			file.Blocks,
		}
		if _, ok := dataBase[file.Name]; !ok {
			dataBase[file.Name] = []Seeder{}
			dataBase[file.Name] = append(dataBase[file.Name], node)
		} else {
			if !Contains(dataBase[file.Name], node) {
				dataBase[file.Name] = append(dataBase[file.Name], node)
			}
		}
	}
}

func Contains(slice []Seeder, seeder Seeder) bool {
	for _, v := range slice {
		if AreSeedersEqual(v, seeder) {
			return true
		}
	}
	return false
}

func AreSeedersEqual(s1 Seeder, s2 Seeder) bool {
	if !s1.Ip.Equal(s2.Ip) {
		return false
	} else if !(s1.Port == s2.Port) {
		return false
	} else if !reflect.DeepEqual(s1.Blocks, s2.Blocks) {
		return false
	}
	return true
}

func RemoveNodeFromList(list []Seeder, ip net.IP) []Seeder {
	var updatedList []Seeder
	for _, item := range list {
		if !item.Ip.Equal(ip) {
			updatedList = append(updatedList, item)
		}
	}
	return updatedList
}
