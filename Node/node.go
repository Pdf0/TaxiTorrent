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

	dirPath, username := GetInitialInfo()

	conn := connectToTracker()
	defer conn.Close()

	SendSyn(conn, dirPath, username)

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
			SendSyn(conn, dirPath, username)
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

func SendSyn(conn net.Conn, dirPath string, username string) {
	syn := CreateSyn(conn, dirPath, username)

	_, err := conn.Write(util.EncodeToBytes(syn))
	checkErr(err)
}

func CreateSyn(conn net.Conn, dirPath string, username string) CentralProtocol.SYN {

	ip, port, nFiles, files := CentralProtocol.GetSYNInfo(conn, dirPath)
	syn := CentralProtocol.CreateSyn(username, ip, port, nFiles, files)

	return syn
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
