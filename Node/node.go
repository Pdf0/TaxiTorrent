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

var SEEDSDIR string
var USERNAME string

func main() {

	SEEDSDIR, USERNAME = GetInitialInfo()

	conn := connectToTracker()
	defer conn.Close()

	SendCentral(conn, "syn")

	clear()
	fmt.Println("Welcome to TaxiTorrent")

	for {
		command := commandLine()

		switch command {
		case "help":
			fmt.Println(" Available commands:\n  help - displays this help menu\n  get - get a file. \n  update - updates your available seeds.\n  clear - clears the screen\n  exit - exits the program")
		case "get":
			var file string
			fmt.Print("File: ")
			fmt.Scanf("%s", &file)
		case "update":
			SendCentral(conn, "update")
		case "clear":
			clear()
		case "exit":
			os.Exit(0)
		}
	}
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

// Função muito javarda mas assim funciona
func SendCentral(conn net.Conn, packetType string) {

	if packetType == "syn" {
		syn := CreateSyn(conn)
		packet := CentralProtocol.CreateCentral("syn", util.EncodeToBytes(syn))
		_, err := conn.Write(util.EncodeToBytes(packet))
		checkErr(err)
	} else if packetType == "update" {
		update := CreateUpdate(conn)
		packet := CentralProtocol.CreateCentral("update", util.EncodeToBytes(update))
		_, err := conn.Write(util.EncodeToBytes(packet))
		checkErr(err)
	}
}

// Estas duas funções podem muito bem fundir-se, assim como as do centralProtocol.go
func CreateSyn(conn net.Conn) CentralProtocol.SYN {

	ip, port, nFiles, files := CentralProtocol.GetNodeInfo(conn, SEEDSDIR)
	syn := CentralProtocol.CreateSyn(USERNAME, ip, port, nFiles, files)

	return syn
}

func CreateUpdate(conn net.Conn) CentralProtocol.Update {
	_, _, nFiles, files := CentralProtocol.GetNodeInfo(conn, SEEDSDIR)
	update := CentralProtocol.CreateUpdate(nFiles, files)

	return update
}

func GetInitialInfo() (string, string) {

	return GetDirPath(), GetUsername()
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
