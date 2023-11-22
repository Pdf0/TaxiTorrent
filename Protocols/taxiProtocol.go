package Protocols

import (
	"net"
	"github.com/google/uuid"
)

/*  id list:
- 0 - syn
- 1 - ack
- 2 - request
- 3 - data
- 4 - maddy (missing)
*/

type TaxiProtocol struct {
	ConnId uuid.UUID
	Id uint8
	Payload []byte
}

type Syn struct {
	Ip net.IP
	FileName string
}

type Request struct {
	Blocklist []uint16
}

type Data struct {
	BlockId uint16
	Block []byte
	Hash string
}

type Maddy struct {
	BlockId uint16
}
