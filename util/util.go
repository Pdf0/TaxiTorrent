package util

import (
	"log"
	"io"
	"bytes"
	"compress/gzip"
	"encoding/gob"
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
    if err != nil {
        log.Fatal(err)
        return err
    }
    return nil
}