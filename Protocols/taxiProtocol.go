package Protocols

import (
	"net"
)

/*  id list:
- 0 - syn
- 1 - ack
- 2 - request
- 3 - data
- 4 - finished
*/

type UDPConnectionInfo struct {
    LocalAddr  net.UDPAddr
    RemoteAddr net.UDPAddr
}

type TaxiProtocol struct {
    SenderIp string
    Id       uint8
    Payload  []byte
}

type SynGates struct {
	Ip       net.IP
	FileName string
}

type Request struct {
	Filename string
	BlocksBF []bool
}

type Data struct {
	Filename string
	BlockId int
	Block   []byte
	Hash    string
}

func CreateSynGates(ip net.IP, fileName string) SynGates {
	return SynGates{
		ip,
		fileName,
	}
}