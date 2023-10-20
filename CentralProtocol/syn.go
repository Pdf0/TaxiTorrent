package CentralProtocol

type SYN struct {
	Username string
	Addr    int64
	Port int
	NFicheiros int
	//struct de ficheiros em chunks
}

func CreateSyn (user string, addr int64,
			port int, nFicheiros int)(syn SYN){
	return SYN{
		user,
		addr,
		port,
		nFicheiros,
	}
}