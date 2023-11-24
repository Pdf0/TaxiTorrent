package main

import (
	"TaxiTorrent/Protocols"
	"TaxiTorrent/util"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	CLIENT_HOST    = "localhost" // dar setup
	CLIENT_TCPPORT = "10001"     // definir um standard

	CLIENT_TYPE = "tcp" // necessario ?

	SERVER_HOST = "10.4.4.2"
	SERVER_PORT = "24"

	CLIENT_UDPPORT = "106" // standard node udp port

	BLOCKSIZE = 1024
)

var SEEDSDIR string
var USERNAME string

func main() {

	//fmt.Println(CLIENT_UDPPORT)

	if len(os.Args) == 3 {
		if !util.DirExists(os.Args[1]) {
			fmt.Println("Non existent Directory")
		} else {

			SEEDSDIR = os.Args[1]
			USERNAME = os.Args[2]

			conn := connectToTracker()
			defer conn.Close()

			SendCentral(conn, "syn")

			clear()
			fmt.Println("Welcome to TaxiTorrent")

			go Listen()

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

						palavras := strings.Fields(command)

						if len(palavras) == 2 {
							//SendCentral(conn, command)

							args := strings.Fields(command)

							// Check if args[1] exists, for example: "> get "
							file := args[1]

							gRequest := Protocols.GetRequest{FileName: file}
							gResponse := new(Protocols.GetResponse)
							commsListandGet(conn, "getrequest", gRequest, gResponse)

							udpconn := connectToSeeder()

							fmt.Println("Conectou-se a um seeder")

							defer udpconn.Close()
							/*
								// Começar a conexão udp com os seeders

								//comecar a gerigonca toda das conversas

								//enviar um syn
								fstmsg := Protocols.CreateSynGates(net.IP(CLIENT_HOST), file)
								fmt.Println("Enviado -> ", fstmsg)
							*/

						} else {
							fmt.Println("Please Specify an argument")
							fmt.Println("> get [file]")
						}

					} else {
						fmt.Println("Invalid command, try using \"help\" to see the available commands")
					}
				}
			}
		}
	} else {
		fmt.Println("The program works as following:")
		fmt.Println("./node [seeds folder] [username]")
	}
}

func connectToTracker() net.Conn {

	conn, err := net.Dial(CLIENT_TYPE, SERVER_HOST+":"+SERVER_PORT)

	util.CheckErr(err)

	return conn
}

func connectToSeeder() net.Conn {
	conn, err := net.Dial("udp", CLIENT_UDPPORT)

	util.CheckErr(err)

	return conn
}

func SendCentral(conn net.Conn, packetType string) {

	if packetType == "syn" {
		syn := CreateSyn(conn)
		packet := Protocols.CreateCentral(packetType, util.EncodeToBytes(syn))
		_, err := conn.Write(util.EncodeToBytes(packet))
		util.CheckErr(err)

	} else if packetType == "update" {
		update := CreateUpdate(conn)
		packet := Protocols.CreateCentral(packetType, util.EncodeToBytes(update))
		_, err := conn.Write(util.EncodeToBytes(packet))
		util.CheckErr(err)

	} else if packetType == "list" {

		lRequest := Protocols.CreateCentral(packetType, []byte{})
		lResponse := new(Protocols.ListResponse)
		commsListandGet(conn, packetType, lRequest, lResponse)

	} else if strings.HasPrefix(packetType, "get") {

		args := strings.Fields(packetType)

		// Check if args[1] exists, for example: "> get "
		file := args[1]

		gRequest := Protocols.GetRequest{FileName: file}
		gResponse := new(Protocols.GetResponse)
		commsListandGet(conn, "getrequest", gRequest, gResponse)

	}
}

func commsListandGet(conn net.Conn, requestType string, requestData interface{}, responseType interface{}) {
	packet := Protocols.CreateCentral(requestType, util.EncodeToBytes(requestData))
	_, err := conn.Write(util.EncodeToBytes(packet))
	util.CheckErr(err)

	buffer := make([]byte, 1024)
	mLen, _ := conn.Read(buffer)

	g := new(Protocols.Central)
	util.DecodeToStruct(buffer[:mLen], g)

	if err := util.DecodeToStruct(g.Payload, responseType); err != nil {
		fmt.Printf("Error decoding %T packet: %s\n", responseType, err.Error())
	}

	fmt.Println(responseType)

}

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

func Listen() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":"+CLIENT_UDPPORT)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}

		packet := buffer[:n]
		fmt.Println("Received data from ", clientAddr, packet)

		//conId, id, payload := handleUPDPpacket(packet)

		//enviar coisas para o outro node

	}
}

func handleUPDPpacket(packet []byte) (uuid.UUID, uint8, []byte) {
	t := new(Protocols.TaxiProtocol)
	util.DecodeToStruct(packet, t)

	return t.ConnId, t.Id, t.Payload
}
