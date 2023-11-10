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

type Central struct {
	PacketType string
	Payload    []byte
}

type File struct {
	Name string
	Size int64 // in bytes
}

func CreateCentral(packetType string, payload []byte) Central {
	return Central{
		packetType,
		payload,
	}
}

func CreateSyn(user string, addr net.IP,
	port uint, nFicheiros int, fileList []File) SYN {
	return SYN{
		user,
		addr,
		port,
		nFicheiros,
		fileList,
	}
}

func ReceiveSYN(syn SYN) bool {
	return true
}

func GetSYNInfo(conn net.Conn, dirPath string) (net.IP, uint, int, []File) {
	ip := util.GetTCPIP(conn)
	port := util.GetTCPPort(conn)

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
