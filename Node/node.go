package main

import (
	"TaxiTorrent/CentralProtocol"
	"TaxiTorrent/util"
	"fmt"
	"net"
	"os"
)

const (
	CLIENT_HOST = "localhost"
	CLIENT_PORT = "10001"
	CLIENT_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "10000"
)

func main() {

	dirPath := GetDirPath()
	username := GetUsername()

	conn := connectToTracker()
	defer conn.Close()

	syn := CreateSyn(conn, dirPath, username)
	central := CentralProtocol.CreateCentral("syn", util.EncodeToBytes(syn))

	_, err := conn.Write(util.EncodeToBytes(central))

	checkErr(err)
}

func connectToTracker() net.Conn {

	conn, err := net.Dial(CLIENT_TYPE, SERVER_HOST+":"+SERVER_PORT)

	checkErr(err)

	return conn
}

func checkErr(err error) {

	if err != nil {

		panic(err)
	}
}

func CreateSyn(conn net.Conn, dirPath string, username string) CentralProtocol.SYN {

	ip, port, nFiles, files := CentralProtocol.GetSYNInfo(conn, dirPath)
	syn := CentralProtocol.CreateSyn(username, ip, port, nFiles, files)

	return syn
}

func GetDirPath() string {
	var path string
	err := false

	for !err {

		fmt.Print("Seeds Directory: ")
		fmt.Scanf("%s", &path)

		err = DirExists("Node/" + path)

	}

	return "Node/" + path
}

func DirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetUsername() string {
	var username string
	fmt.Print("Username: ")
	fmt.Scanf("%s", &username)

	return username
}
