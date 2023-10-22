package main

import (
	"TaxiTorrent/CentralProtocol"
	"TaxiTorrent/util"
	"net"
)

const (
	CLIENT_HOST = "localhost"
	CLIENT_PORT = "10001"
	CLIENT_TYPE = "tcp"
	SERVER_HOST = "localhost"
	SERVER_PORT = "10000"
)

func main() {
	/*
		var buffer string

			for {
				fmt.Scanln(&buffer)

				con, err := net.Dial("tcp", SERVER_HOST+":"+SERVER_PORT)

				checkErr(err)

				defer con.Close()

				_, err = con.Write([]byte(buffer))

				checkErr(err)
			}
	*/

	syn := CentralProtocol.CreateSyn("gato", 2, 22, 20, []CentralProtocol.File{{Name: "Pedro", Size: 3}})

	con, err := net.Dial("tcp", SERVER_HOST+":"+SERVER_PORT)

	checkErr(err)

	defer con.Close()

	_, err = con.Write(util.EncodeToBytes(syn))

	checkErr(err)
}

func checkErr(err error) {

	if err != nil {

		panic(err)
	}
}
