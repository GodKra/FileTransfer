package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/jhoonb/archivex"
	"io"
	"log"
	"net"
	"os"
)

var (
	path   = flag.String("path", "FileTransfer", "Usage : -path <path> eg : -path FileTransfer/test")
	ipFlag = flag.String("ip", "localhost:5555", "Usage : -ip <ip address:port> eg : -ip localhost:5555")
)

func main() {
	flag.Parse()

	s, e := os.Stat(*path)
	checkError(e)
	ip := *ipFlag
	name := "temp"

	var z zipper
	if s.IsDir() {
		z = dir(*path)
	} else {
		z = file(*path)
	}
	z.zip(name, *path)

	f, e := os.Open(name + ".zip")
	defer f.Close()
	checkError(e)

	defer os.Remove(name + ".zip")

	conn := connect(ip)
	stat, e := f.Stat()
	checkError(e)

	e = sendFileSize(stat.Size(), conn)
	checkError(e)

	fmt.Println(stat.Size())
	i, e := io.Copy(conn, f)
	checkError(e)

	fmt.Println("Sending...")

	fmt.Printf("Succesfully sent %v bytes of data\n", i)
	conn.Close()
}

func connect(ip string) net.Conn {
	conn, e := net.Dial("tcp", ip)
	checkError(e)
	return conn
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func sendFileSize(size int64, conn net.Conn) error {
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], uint64(size))
	_, e := conn.Write(b[:])
	return e
}

type zipper interface {
	zip(name, path string)
}

type dir string

func (d dir) zip(name, path string) {
	zipfile := archivex.ZipFile{}
	e := zipfile.Create(name)
	checkError(e)
	e = zipfile.AddAll(path, false)
	checkError(e)
	zipfile.Close()
}

type file string

func (f file) zip(name string, path string) {
	zipfile := archivex.ZipFile{}
	e := zipfile.Create(name)
	checkError(e)
	e = zipfile.AddFile(path)
	checkError(e)
	zipfile.Close()
}