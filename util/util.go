package util

import (
	"fmt"
	"log"
	"io"
	"bytes"
	"compress/gzip"
	"encoding/gob"
)

func EncodeToBytes(p interface{}) []byte {

    buf := bytes.Buffer{}
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(p)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
    return Compress(buf.Bytes())
}

func Compress(s []byte) []byte {

    zipbuf := bytes.Buffer{}
    zipped := gzip.NewWriter(&zipbuf)
    zipped.Write(s)
    zipped.Close()
    fmt.Println("compressed size (bytes): ", len(zipbuf.Bytes()))
    return zipbuf.Bytes()
}

func Decompress(s []byte) []byte {

    rdr, _ := gzip.NewReader(bytes.NewReader(s))
    data, err := io.ReadAll(rdr)
    if err != nil {
        log.Fatal(err)
    }
    rdr.Close()
    fmt.Println("uncompressed size (bytes): ", len(data))
    return data
}

func DecodeToStruct(s []byte, i interface{}) error {
    dec := gob.NewDecoder(bytes.NewReader(Decompress(s)))
    err := dec.Decode(i)
    if err != nil {
        log.Fatal(err)
        return err
    }
    return nil
}