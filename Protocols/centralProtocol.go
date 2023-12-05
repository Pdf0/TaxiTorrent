package Protocols

import (
	"TaxiTorrent/util"
	"log"
	"math"
	"net"
	"os"
)

const (
	BLOCKSIZE = 1024
)

type FileInfo struct {
	FileSize    uint64
	SeedersInfo []Seeder
}

type Seeder struct {
	Ip              net.IP
	Username        string
	Port            uint
	BlocksAvailable []bool
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

type BlockUpdate struct {
	Filename string
	BlockId  int
}

type Central struct {
	PacketType string
	Payload    []byte
}

type File struct {
	Name    string
	Size    int64
	NBlocks int64
	Blocks  []bool
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
			GetBlocksBitfield(dirPath + "/" + file.Name()),
		}
		fileCount++
	}
	return ip, port, fileCount, filesArray
}

func GetFileNBlocks(fileSize int64) uint64 {
	return uint64(math.Ceil((float64(fileSize)) / float64(BLOCKSIZE)))
}

func GetBlocksBitfield(fp string) []bool {
	data, _ := os.ReadFile(fp)
	blocks := make([]bool, int(math.Ceil(float64(len(data))/float64(BLOCKSIZE))))
	for i := 0; i < len(blocks); i += 1 {
		blocks[i] = true
	}
	return blocks
}

func DeepCopySeeders(seeders []Seeder) []Seeder {
	newSeeder := make([]Seeder, len(seeders))

	for i, seeder := range seeders {
		newSeeder[i] = DeepCopySeeder(seeder)
	}
	return newSeeder
}

func DeepCopySeeder(s Seeder) Seeder {
	return Seeder{Ip: s.Ip, Port: s.Port, BlocksAvailable: s.BlocksAvailable}
}

func QueryUsername(username string) []string {
	ips, err := net.LookupHost(username)
	if err != nil {
		return nil
	}
	return ips
}

func QueryIp(ip net.IP) string {
	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return ""
	}

	if len(names) > 0 {
		return names[0]
	}

	return ""
}
