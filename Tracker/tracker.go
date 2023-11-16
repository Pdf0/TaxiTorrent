package main

import (
	"TaxiTorrent/Protocols"
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

func main() {

	dataBase := make(map[string][]Protocols.Seeder)

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

func processClient(connection net.Conn, dataBase map[string][]Protocols.Seeder) {
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

		g := new(Protocols.Central)
		util.DecodeToStruct(buffer[:mLen], g)

		switch g.PacketType {
		case "syn":
			s := new(Protocols.SYN)
			if err := util.DecodeToStruct(g.Payload, s); err != nil {
				fmt.Println("Error decoding SYN packet:", err.Error())
				continue
			}
			InsertNodeInDataBase(dataBase, s.FileList, s.Ip, s.Port)
			fmt.Println("Received: ", *s)
			fmt.Println(dataBase)

		case "update":
			u := new(Protocols.Update)
			if err := util.DecodeToStruct(g.Payload, u); err != nil {
				fmt.Println("Error decoding Update packet:", err.Error())
				continue
			}
			UpdateNode(dataBase, u.FileList, ip, port)
			fmt.Println("Updated Node (" + fullAddr + ")")
			fmt.Println(dataBase)

		case "list":
			fmt.Println("ListRequest from ", fullAddr)
			var fileList []string
			for fileName := range dataBase {
				fileList = append(fileList, fileName)
			}
			lResponse := Protocols.ListResponse{FileList: fileList}
			central := Protocols.Central{PacketType: "ListResponse", Payload: util.EncodeToBytes(lResponse)}
			_, err := connection.Write(util.EncodeToBytes(central))
			util.CheckErr(err)
			fmt.Println("ListResponse to ", fullAddr, ": ", lResponse)

		case "getrequest":
			gRequest := new(Protocols.GetRequest)
			if err := util.DecodeToStruct(g.Payload, gRequest); err != nil {
				fmt.Println("Error decoding GetRequest packet:", err.Error())
				continue
			}
			fmt.Println("GetRequest from ", fullAddr, ": ", gRequest)
			var seedersList []Protocols.Seeder
			for fileName, seeders := range dataBase {
				if fileName == gRequest.FileName {
					seedersList = append(seedersList, seeders...)
				}
			}
			gResponse := Protocols.GetResponse{Seeders: seedersList}
			central := Protocols.Central{PacketType: "GetResponse", Payload: util.EncodeToBytes(gResponse)}
			_, err := connection.Write(util.EncodeToBytes(central))
			util.CheckErr(err)
			fmt.Println("GetResponse to ", fullAddr, ": ", gResponse)
		}
	}
}

func UpdateNode(dataBase map[string][]Protocols.Seeder, files []Protocols.File, ip net.IP, port uint) {
	DisconnectNode(dataBase, ip, port)
	InsertNodeInDataBase(dataBase, files, ip, port)
}

func DisconnectNode(dataBase map[string][]Protocols.Seeder, ip net.IP, port uint) {
	for file, list := range dataBase {
		dataBase[file] = RemoveNodeFromList(list, ip)
		if len(dataBase[file]) == 0 {
			delete(dataBase, file)
		}
	}
}

func InsertNodeInDataBase(dataBase map[string][]Protocols.Seeder, files []Protocols.File, ip net.IP, port uint) {
	for _, file := range files {
		node := Protocols.Seeder{
			Ip:     ip,
			Port:   port,
			Blocks: file.Blocks,
		}
		if _, ok := dataBase[file.Name]; !ok {
			dataBase[file.Name] = []Protocols.Seeder{}
			dataBase[file.Name] = append(dataBase[file.Name], node)
		} else {
			if !Contains(dataBase[file.Name], node) {
				dataBase[file.Name] = append(dataBase[file.Name], node)
			}
		}
	}
}

func Contains(slice []Protocols.Seeder, seeder Protocols.Seeder) bool {
	for _, v := range slice {
		if AreSeedersEqual(v, seeder) {
			return true
		}
	}
	return false
}

func AreSeedersEqual(s1 Protocols.Seeder, s2 Protocols.Seeder) bool {
	if !s1.Ip.Equal(s2.Ip) {
		return false
	} else if !(s1.Port == s2.Port) {
		return false
	} else if !reflect.DeepEqual(s1.Blocks, s2.Blocks) {
		return false
	}
	return true
}

func RemoveNodeFromList(list []Protocols.Seeder, ip net.IP) []Protocols.Seeder {
	var updatedList []Protocols.Seeder
	for _, item := range list {
		if !item.Ip.Equal(ip) {
			updatedList = append(updatedList, item)
		}
	}
	return updatedList
}
