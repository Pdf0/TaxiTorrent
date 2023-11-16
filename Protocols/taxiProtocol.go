package Protocols

// keep-alive: <len=0000>
// The keep-alive message is a message with zero bytes, specified with the length prefix set to zero. There is no message ID and no payload. Peers may close a connection if they receive no messages (keep-alive or any other message) for a certain period of time, so a keep-alive message must be sent to maintain the connection alive if no command have been sent for a given amount of time. This amount of time is generally two minutes.

// choke: <len=0001><id=0>
// The choke message is fixed-length and has no payload.

// unchoke: <len=0001><id=1>
// The unchoke message is fixed-length and has no payload.

// interested: <len=0001><id=2>
// The interested message is fixed-length and has no payload.

// not interested: <len=0001><id=3>
// The not interested message is fixed-length and has no payload.

// have: <len=0005><id=4><piece index>
// The have message is fixed length. The payload is the zero-based index of a piece that has just been successfully downloaded and verified via the hash.
// Only use it when receiving pieces so it wants to update all other nodes

// request: <len=0013><id=6><index><begin><length>
// The request message is fixed length, and is used to request a block. The payload contains the following information:
// index: integer specifying the zero-based piece index
// begin: integer specifying the zero-based byte offset within the piece
// length: integer specifying the requested length

// piece: <len=0009+X><id=7><index><begin><block>
// The piece message is variable length, where X is the length of the block. The payload contains the following information:
// index: integer specifying the zero-based piece index
// begin: integer specifying the zero-based byte offset within the piece
// block: block of data, which is a subset of the piece specified by index.

// cancel: <len=0013><id=8><index><begin><length>
// The cancel message is fixed length, and is used to cancel block requests. The payload is identical to that of the "request" message. It is typically used during "End Game"

type taxiProtocol struct {
	length uint8
	id uint8
	payload []byte
}

type handshake struct {
	pstrl uint8
	pstr string
	fileName string
}

type piecesBitfield struct {
	bitfield []bool
}

type request struct {
	index uint8
	begin uint8
	length uint16
}

type cancel struct {
	index uint8
	begin uint8
	length uint8
}

type piece struct {
	index uint8
	begin uint8
	block []byte
}
