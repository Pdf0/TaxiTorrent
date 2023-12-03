package main

import (
	"TaxiTorrent/Protocols"
	"TaxiTorrent/util"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	CLIENT_TYPE = "tcp" // necessario ?

	// valores definidos na topologia para o servidor 2
	SERVER_IP   = "10.4.4.2"
	SERVER_PORT = "9090"

	CLIENT_UDPPORT = "9090"

	BLOCKSIZE = 1024
)

type Connection struct {
	conn             *net.UDPConn
	FileName         string
	FileSize         uint64 // bytes
	BlocksDownloaded []byte
	BlocksToDownload []bool
}

var CLIENT_HOST string
var SEEDSDIR string
var USERNAME string
var ackReceived chan bool
var ackMutex sync.Mutex
var closeOnce sync.Once

func main() {

	CLIENT_HOST = getPublicIP()
	dataBase := createDataBase()

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

							args := strings.Fields(command)

							//FIX: Check if args[1] exists, for example: "> get "
							file := args[1]

							//TODO: Verificar se o Node já tem o ficheiro que está a pedir

							gRequest := Protocols.GetRequest{FileName: file}
							gResponse := new(Protocols.GetResponse)
							commsListandGet(conn, "getrequest", gRequest, gResponse)

							sendInitialSynPackets(gResponse, file, &dataBase)

							//TODO: Verificar se todos os handshakes foram bem sucedidos

							var wg sync.WaitGroup
							for nodeIp, connection := range dataBase {
								wg.Add(1)
								go sendRequestConcurrent(nodeIp, connection, &wg)
							}
							wg.Wait()
							//TODO: Verificar se todos os blocos estão downloaded

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

func getPublicIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func connectToTracker() net.Conn {

	conn, err := net.Dial(CLIENT_TYPE, SERVER_IP+":"+SERVER_PORT)

	util.CheckErr(err)

	return conn
}

func connectToSeeder(udpAddr *net.UDPAddr) *net.UDPConn {

	conn, err := net.DialUDP("udp", nil, udpAddr)

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
		fmt.Println("Received data from ", clientAddr)

		go handleUDPpacket(Protocols.UDPConnectionInfo{LocalAddr: *serverAddr, RemoteAddr: *clientAddr}, packet)
	}
}

func handleUDPpacket(connInfo Protocols.UDPConnectionInfo, packet []byte) {
	t := new(Protocols.TaxiProtocol)
	util.DecodeToStruct(packet, t)

	// Syn
	if t.Id == 0 {

		taxiAck := createAck(connInfo)
		sendPacketOverUDP(connInfo.RemoteAddr, util.EncodeToBytes(taxiAck))

		fmt.Println("Ack sent to", connInfo.RemoteAddr.String())
		// Ack
	} else if t.Id == 1 {
		fmt.Println("ACK received")
		ackMutex.Lock()
		closeOnce.Do(func() {
			close(ackReceived)
		})
		ackMutex.Unlock()
		// Request
	} else if t.Id == 2 {
		request := new(Protocols.Request)
		util.DecodeToStruct(t.Payload, request)

		fmt.Println("Received a Request: ", request)
		handleRequest(connInfo, request)
		// Data
	} else if t.Id == 3 {
		data := new(Protocols.Data)
		util.DecodeToStruct(t.Payload, data)

		writeDataToFile(data)
	}
}

func createAck(connInfo Protocols.UDPConnectionInfo) Protocols.TaxiProtocol {
	return Protocols.TaxiProtocol{
		SenderIp: connInfo.LocalAddr.IP.String(),
		Id:       uint8(1),
	}
}

func createDataBase() map[string]Connection {
	return make(map[string]Connection)
}

func sendInitialSynPackets(gResponse *Protocols.GetResponse, fileName string, dataBase *map[string]Connection) {

	//Algoritmo de distribuição de blocos (Para já só distribui os blocos continua e uniformente)
	blocksPerNode := (int(gResponse.Size) / BLOCKSIZE) / len(gResponse.Seeders)
	blocksOffset := 0

	var wg sync.WaitGroup
	for _, node := range gResponse.Seeders {
		wg.Add(1)
		go makeHandshake(node, fileName, gResponse.Size, blocksPerNode, blocksOffset, dataBase, &wg)
		blocksOffset = blocksOffset + blocksPerNode
	}
	wg.Wait()
}

func makeHandshake(node Protocols.Seeder, fileName string, fileSize uint64, nBlocks int, blocksOffset int, dataBase *map[string]Connection, wg *sync.WaitGroup) {

	defer wg.Done()
	udpAddr, _ := net.ResolveUDPAddr("udp", node.Ip.String()+":"+CLIENT_UDPPORT)
	udpconn := connectToSeeder(udpAddr)

	fmt.Println("P2P connection established with", udpAddr)

	sendInitialSynPacket(fileName, udpconn)

	ackReceived = make(chan bool)

	retries := 3
	ackReceivedFlag := false

retryloop:
	for i := 0; i < retries; i++ {
		select {
		case <-ackReceived:
			fmt.Println("ACK received")
			ackReceivedFlag = true
			break retryloop
		case <-time.After(1 * time.Second):
			fmt.Println("Resending SynGate")
			sendInitialSynPacket(fileName, udpconn)
		}
	}
	if !ackReceivedFlag {
		fmt.Println("Closing connection, no ACK received.")
		udpconn.Close()
		return
	}
	blocksToDownload := util.CreateBitFieldFromTo(node.BlocksAvailable, blocksOffset, blocksOffset+nBlocks)
	(*dataBase)[node.Ip.String()] = Connection{
		udpconn,
		fileName,
		fileSize,
		make([]byte, 0),
		blocksToDownload,
	}
}

func sendInitialSynPacket(fileName string, udpconn *net.UDPConn) *net.UDPConn {

	synGate := Protocols.CreateSynGates(net.IP(CLIENT_HOST), fileName)
	taxi := Protocols.TaxiProtocol{SenderIp: getPublicIP(), Id: 0, Payload: util.EncodeToBytes(synGate)}
	_, err := udpconn.Write(util.EncodeToBytes(taxi))

	util.CheckErr(err)

	fmt.Println("SynGate sent successfully")

	return udpconn
}

func sendRequestConcurrent(nodeIP string, connection Connection, wg *sync.WaitGroup) {
	defer wg.Done()
	sendRequest(nodeIP, connection)
}

func sendRequest(nodeIP string, connection Connection) {

	request := Protocols.Request{Filename: connection.FileName, BlocksBF: connection.BlocksToDownload}
	taxi := Protocols.TaxiProtocol{
		SenderIp: getPublicIP(),
		Id:       2,
		Payload:  util.EncodeToBytes(request),
	}
	_, err := connection.conn.Write(util.EncodeToBytes(taxi))

	util.CheckErr(err)
}

func handleRequest(connInfo Protocols.UDPConnectionInfo, request *Protocols.Request) {

	for i, value := range request.BlocksBF {
		if value {
			data := createDataPacket(request.Filename, i)
			sendBlock(data, connInfo)
		}
	}
}

func sendBlock(data Protocols.Data, connInfo Protocols.UDPConnectionInfo) {
	taxi := Protocols.TaxiProtocol{
		SenderIp: connInfo.LocalAddr.String(),
		Id:       3,
		Payload:  util.EncodeToBytes(data),
	}
	sendPacketOverUDP(connInfo.RemoteAddr, util.EncodeToBytes(taxi))
}

func createDataPacket(file string, blockId int) Protocols.Data {

	startIndex := int64(blockId * BLOCKSIZE)
	fp := fmt.Sprintf("%s/%s", SEEDSDIR, file)

	fileData, _ := os.Open(fp)

	defer fileData.Close()

	fileData.Seek(startIndex, 0)
	buffer := make([]byte, BLOCKSIZE)
	fileData.Read(buffer)

	return Protocols.Data{
		Filename: file,
		BlockId:  blockId,
		Block:    buffer,
		Hash:     util.HashBlockMD5(buffer),
	}
}

func sendPacketOverUDP(addr net.UDPAddr, data []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr.IP.String()+":"+CLIENT_UDPPORT)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func writeDataToFile(data *Protocols.Data) {

	download_path := fmt.Sprintf("%s/%s", SEEDSDIR, data.Filename)
	file, _ := os.Create(download_path)

	defer file.Close()

	writeBlockToFile(file, (*data).Block, int64(data.BlockId)*BLOCKSIZE)

}

func writeBlockToFile(file *os.File, block []byte, offset int64) {
	file.Seek(offset, io.SeekStart)
	file.Write(block)
}
