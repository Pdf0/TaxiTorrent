package Protocols

import (
	"TaxiTorrent/util"
	"log"
	"math"
	"net"
	"os"
)

const (
	BLOCKSIZE = 256
)

type FileInfo struct {
	FileSize uint64
	SeedersInfo []Seeder
}

type Seeder struct {
	Ip     net.IP
	Port   uint
	BlocksAvailable []string
	BlocksToDownload []bool
}

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

type GetRequest struct {
	FileName string
}

type GetResponse struct {
	Seeders []Seeder
	Size    uint64
}

type ListResponse struct {
	FileList []string
}

type Central struct {
	PacketType string
	Payload    []byte
}

type File struct {
	Name    string
	Size    int64
	NBlocks int64
	Blocks  []string
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

func CreateEmptyCentral() Central { return Central{} }

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
	ip, port := util.GetTCPLocalAddr(conn)

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
			file.Name(),
			fileInfo.Size(),
			int64(GetFileNBlocks(fileInfo.Size())),
			GetBlocksHashes(dirPath + "/" + file.Name()),
		}
		fileCount++
	}

	return ip, port, fileCount, filesArray
}

func GetFileNBlocks(fileSize int64) uint64 {
	return uint64(math.Ceil((float64(fileSize)) / float64(BLOCKSIZE)))
}

func GetBlocksHashes(fp string) []string {
	data, _ := os.ReadFile(fp)

	var blocks []string

	for i := 0; i < len(data); i += int(BLOCKSIZE) {
		end := i + int(BLOCKSIZE)
		if end > len(data) {
			end = len(data)
		}
		block := data[i:end]
		blocks = append(blocks, util.HashBlockMD5(block))
	}
	return blocks
}
