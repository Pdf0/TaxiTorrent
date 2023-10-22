package CentralProtocol

type SYN struct {
	Username string
	Addr    int64
	Port int
	NFicheiros int
	FileList[] File
}

type File struct {
	Name string
	Size int // in bytes
}

type ACK struct {
	Error error
}

func CreateSyn (user string, addr int64,
			port int, nFicheiros int, fileList []File) SYN {
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