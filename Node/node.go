package main

import (
	"TaxiTorrent/CentralProtocol"
	"TaxiTorrent/util"
	"fmt"
	"net"
)

const (
	CLIENT_HOST = "localhost"
	CLIENT_PORT = "10001"
	CLIENT_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "10000"
)

func main() {
	conn := connectToTracker()

	defer conn.Close()

	syn := CreateSyn(conn)

	_, err := conn.Write(util.EncodeToBytes(syn))

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

func CreateSyn(conn net.Conn) CentralProtocol.SYN {

	dirPath := "Node/"
	var tmp string
	fmt.Print("Seeds Directory: ")
	fmt.Scanf("%s", &tmp)
	dirPath = dirPath + tmp
	name, ip, port, nFiles, files := CentralProtocol.GetSYNInfo(conn, dirPath)
	syn := CentralProtocol.CreateSyn(name, ip, port, nFiles, files)

	return syn
}
