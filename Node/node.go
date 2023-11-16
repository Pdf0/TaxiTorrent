package main

import (
	"TaxiTorrent/Protocols"
	"TaxiTorrent/util"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CLIENT_HOST = "localhost"
	CLIENT_PORT = "10001"
	CLIENT_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "10000"
	BLOCKSIZE   = 1024
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
		case "help\n":
			fmt.Println(" Available commands:\n  help - displays this help menu\n  get <file> - get a file. Example: get file1.txt. \n  update - updates your available seeds.\n  clear - clears the screen\n  exit - exits the program")
		case "update\n":
			SendCentral(conn, "update")
		case "list\n":
			SendCentral(conn, "list")
		case "clear\n":
			clear()
		case "exit\n":
			os.Exit(0)
		default:
			if strings.HasPrefix(command, "get") {
				SendCentral(conn, command)

				// Começar a conexão udp com os seeders
			} else {
				fmt.Println("Invalid command, try using \"help\" to see the available commands")
			}
		}
	}
}

func connectToTracker() net.Conn {

	conn, err := net.Dial(CLIENT_TYPE, SERVER_HOST+":"+SERVER_PORT)

	util.CheckErr(err)

	return conn
}

// Função muito javarda mas assim funciona
func SendCentral(conn net.Conn, packetType string) {

	if packetType == "syn" {
		syn := CreateSyn(conn)
		packet := Protocols.CreateCentral("syn", util.EncodeToBytes(syn))
		_, err := conn.Write(util.EncodeToBytes(packet))
		util.CheckErr(err)

	} else if packetType == "update" {
		update := CreateUpdate(conn)
		packet := Protocols.CreateCentral("update", util.EncodeToBytes(update))
		_, err := conn.Write(util.EncodeToBytes(packet))
		util.CheckErr(err)

	} else if packetType == "list" {
		packet := Protocols.CreateCentral("list", []byte{})
		_, err := conn.Write(util.EncodeToBytes(packet))
		util.CheckErr(err)

		// mover esta função para fora daqui, sai bicho
		// fazer uma função para isto, tal como se repete no tracker
		buffer := make([]byte, 1024)
		mLen, _ := conn.Read(buffer)

		g := new(Protocols.Central)
		util.DecodeToStruct(buffer[:mLen], g)
		lResponse := new(Protocols.ListResponse)
		if err := util.DecodeToStruct(g.Payload, lResponse); err != nil {
			fmt.Println("Error decoding ListResponse packet:", err.Error())
		}
		fmt.Println(*lResponse)

	} else if strings.HasPrefix(packetType, "get") {
		args := strings.Fields(packetType)
		// Checkar se args[1] realmente existe. Ex: "> get "
		file := args[1]
		packet := Protocols.CreateCentral("getrequest", util.EncodeToBytes(Protocols.GetRequest{FileName: file}))
		_, err := conn.Write(util.EncodeToBytes(packet))
		util.CheckErr(err)

		// mover esta função para fora daqui, sai bicho
		// fazer uma função para isto, tal como se repete no tracker
		buffer := make([]byte, 1024)
		mLen, _ := conn.Read(buffer)

		g := new(Protocols.Central)
		util.DecodeToStruct(buffer[:mLen], g)
		gResponse := new(Protocols.GetResponse)
		if err := util.DecodeToStruct(g.Payload, gResponse); err != nil {
			fmt.Println("Error decoding GetResponse packet:", err.Error())
		}
		fmt.Println(*gResponse)
	}
}

// Estas duas funções podem muito bem fundir-se, assim como as do Protocols.go
func CreateSyn(conn net.Conn) Protocols.SYN {

	ip, port, nFiles, files := Protocols.GetNodeInfo(conn, SEEDSDIR)
	syn := Protocols.CreateSyn(USERNAME, ip, port, nFiles, files)

	return syn
}

func CreateUpdate(conn net.Conn) Protocols.Update {
	_, _, nFiles, files := Protocols.GetNodeInfo(conn, SEEDSDIR)
	update := Protocols.CreateUpdate(nFiles, files)

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
