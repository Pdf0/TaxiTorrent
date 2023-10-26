package CentralProtocol

import (
	"TaxiTorrent/util"
	"fmt"
	"log"
	"net"
	"os"
)

type SYN struct {
	Username string
	Ip    net.IP
	Port uint
	NFicheiros int
	FileList[] File
}

type File struct {
	Name string
	Size int64 // in bytes
}

type ACK struct {
	Error error
}

func CreateSyn (user string, addr net.IP,
			port uint, nFicheiros int, fileList []File) SYN {
	return SYN{
		user,
		addr,
		port,
		nFicheiros,
		fileList,
	}
}

func ReceiveSYN (syn SYN) bool {
	return true
}

func GetSYNInfo (conn net.Conn, dirPath string) (string, net.IP, uint, int, []File){
	var name string
	
	fmt.Print("Username: ")
	fmt.Scanf("%s", &name)
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
		fileInfo, _  := file.Info()
		filesArray[fileCount] = File{
			Name: file.Name(),
			Size: fileInfo.Size(),
		}
		fileCount++
	}

	return name, ip, port, fileCount, filesArray
}