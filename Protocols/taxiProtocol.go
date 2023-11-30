package Protocols

import (
	"net"
)

/*  id list:
- 0 - syn
- 1 - ack
- 2 - request
- 3 - data
- 4 - maddy (missing)
*/

type UDPConnectionInfo struct {
    LocalAddr  string
    RemoteAddr string
}

type TaxiProtocol struct {
    ConnInfo UDPConnectionInfo
    Id       uint8
    Payload  []byte
}

type SynGates struct {
	Ip       net.IP
	FileName string
}

type Request struct {
	Blocklist []int
}

type Data struct {
	BlockId int
	Block   []byte
	Hash    string
}

type Maddy struct {
	BlockId uint16
}

func CreateSynGates(ip net.IP, fileName string) SynGates {
	return SynGates{
		ip,
		fileName,
	}
}

func GetUDPConnInfo(conn *net.UDPConn) UDPConnectionInfo {
	return UDPConnectionInfo{
		LocalAddr:  conn.LocalAddr().String(),
		RemoteAddr: conn.RemoteAddr().String(),
	}
}