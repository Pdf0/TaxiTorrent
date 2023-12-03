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
	SERVER_HOST = "10.4.4.2" // servidor2
	SERVER_PORT = "9090"
	SERVER_TYPE = "tcp"
	BLOCKSIZE   = 256
)

func main() {

	dataBase := make(map[string]*Protocols.FileInfo)

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
		go processClient(connection, &dataBase)
	}
}

func processClient(connection net.Conn, dataBase *map[string]*Protocols.FileInfo) {
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
			fmt.Println(*dataBase)

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
			for fileName := range *dataBase {
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
			var fileSize uint64
			for fileName, fileInfo := range *dataBase {
				if fileName == gRequest.FileName {
					seedersList = Protocols.DeepCopySeeders(fileInfo.SeedersInfo)
					fileSize = fileInfo.FileSize
					fmt.Println(fileInfo.SeedersInfo)
				}
			}
			gResponse := Protocols.GetResponse{Seeders: seedersList, Size: fileSize}
			central := Protocols.Central{PacketType: "GetResponse", Payload: util.EncodeToBytes(gResponse)}
			_, err := connection.Write(util.EncodeToBytes(central))
			util.CheckErr(err)
			fmt.Println("GetResponse to ", fullAddr, ": ", gResponse)
		}
	}
}

func UpdateNode(dataBase *map[string]*Protocols.FileInfo, files []Protocols.File, ip net.IP, port uint) {
	DisconnectNode(dataBase, ip, port)
	InsertNodeInDataBase(dataBase, files, ip, port)
}

func DisconnectNode(dataBase *map[string]*Protocols.FileInfo, ip net.IP, port uint) {
	for file, fileInfo := range *dataBase {
		(*dataBase)[file].SeedersInfo = RemoveNodeFromList(fileInfo.SeedersInfo, ip)
		if len((*dataBase)[file].SeedersInfo) == 0 {
			delete(*dataBase, file)
		}
	}
}

func InsertNodeInDataBase(dataBase *map[string]*Protocols.FileInfo, files []Protocols.File, ip net.IP, port uint) {
	for _, file := range files {
		node := Protocols.Seeder{
			Ip:              ip,
			Port:            port,
			BlocksAvailable: file.Blocks,
		}
		if _, ok := (*dataBase)[file.Name]; !ok {
			(*dataBase)[file.Name] = &Protocols.FileInfo{FileSize: uint64(file.Size), SeedersInfo: []Protocols.Seeder{node}}
		} else {
			if !Contains((*dataBase)[file.Name].SeedersInfo, node) {
				(*dataBase)[file.Name].SeedersInfo = append((*dataBase)[file.Name].SeedersInfo, node)
			}
		}
		fmt.Println((*dataBase)[file.Name])
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
	} else if !reflect.DeepEqual(s1.BlocksAvailable, s2.BlocksAvailable) {
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
