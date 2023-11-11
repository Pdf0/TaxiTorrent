package CentralProtocol

import (
	"TaxiTorrent/util"
	"log"
	"net"
	"os"
)

type SYN struct {
	Username   string
	Ip         net.IP
	Port       uint
	NFicheiros int
	FileList   []File
}

type Update struct {
	NFicheiros int
	FileList   []File
}

type Central struct {
	PacketType string
	Payload    []byte
}

type File struct {
	Name string
	Size int64 // in bytes
}

func FillCentral(central Central, packetType string, payload []byte) {
	central.PacketType = packetType
	central.Payload = payload
}

func CreateCentral(packetType string, payload []byte) Central {
	return Central{
		packetType,
		payload,
	}
}

func CreateEmptyCentral() Central {return Central{}}

func CreateSyn(user string, ip net.IP,
	port uint, nFicheiros int, fileList []File) SYN {
	return SYN{
		user,
		ip,
		port,
		nFicheiros,
		fileList,
	}
}

func CreateUpdate(nFicheiros int, fileList []File) Update {
	return Update{
		nFicheiros,
		fileList,
	}
}

func GetNodeInfo(conn net.Conn, dirPath string) (net.IP, uint, int, []File) {
	ip, port := util.GetTCPAddr(conn)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	fileCount := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileCount++
	}

	filesArray := make([]File, fileCount)
	fileCount = 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileInfo, _ := file.Info()
		filesArray[fileCount] = File{
			Name: file.Name(),
			Size: fileInfo.Size(),
		}
		fileCount++
	}

	return ip, port, fileCount, filesArray
}
