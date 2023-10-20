package CentralProtocol

type SYN struct {
	username string
	addr    int64
	port int
	nFicheiros int
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