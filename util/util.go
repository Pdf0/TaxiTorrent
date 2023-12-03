package util

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func EncodeToBytes(i interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	if err != nil {
		log.Fatal(err)
	}
	return Compress(buf.Bytes())
}

func Compress(s []byte) []byte {

	zipbuf := bytes.Buffer{}
	zipped := gzip.NewWriter(&zipbuf)
	zipped.Write(s)
	zipped.Close()
	return zipbuf.Bytes()
}

func Decompress(s []byte) []byte {

	rdr, _ := gzip.NewReader(bytes.NewReader(s))
	data, err := io.ReadAll(rdr)
	if err != nil {
		log.Fatal(err)
	}
	rdr.Close()
	return data
}

func DecodeToStruct(s []byte, i interface{}) error {
	dec := gob.NewDecoder(bytes.NewReader(Decompress(s)))
	err := dec.Decode(i)
	if err != nil && err != io.EOF {
		fmt.Println("Error decoding:", err.Error())
		return err
	}
	return nil
}

func GetTCPLocalAddr(conn net.Conn) (net.IP, uint) {
	return GetTCPLocalIP(conn), GetTCPLocalPort(conn)
}

func GetTCPLocalIP(conn net.Conn) net.IP {
	return conn.LocalAddr().(*net.TCPAddr).IP
}

func GetTCPLocalPort(conn net.Conn) uint {
	return uint(conn.LocalAddr().(*net.TCPAddr).Port)
}

func GetTCPRemoteAddr(conn net.Conn) (net.IP, uint) {
	return GetTCPRemoteIP(conn), GetTCPRemotePort(conn)
}

func GetTCPRemoteIP(conn net.Conn) net.IP {
	return conn.RemoteAddr().(*net.TCPAddr).IP
}

func GetTCPRemotePort(conn net.Conn) uint {
	return uint(conn.RemoteAddr().(*net.TCPAddr).Port)
}

func HashBlockMD5(block []byte) string {
	hasher := md5.New()
	hasher.Write(block)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func DirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CreateBitFieldFromTo(bitField []bool, from int, to int) []bool {
	for i := 0; i <= len(bitField); i++ {
		if from > i && i <= to {
			bitField[i] = false
		}
	}
	newBitField := make([]bool, len(bitField))
	copy(newBitField, bitField)

	return newBitField
}