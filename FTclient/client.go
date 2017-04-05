package main

import (
	"archive/zip"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

var (
	path      = flag.String("path", "FileTransfer", "Usage : -path <path> | eg : -path FileTransfer/test")
	ipFlag    = flag.String("ip", "localhost:5555", "Usage : -ip <ip address:port> | eg : -ip localhost:5555")
)

const name = "temp.zip"

func main() {
	flag.Parse()
	ip := *ipFlag
	conn := connect(ip)
	defer conn.Close()

	fmt.Println("Zipping File..")
	f, e := archive(name, *path)
	checkError(e)
	defer os.Remove(name)
	defer f.Close()
	_, e = f.Seek(0, io.SeekStart)
	checkError(e)
	stat, e := f.Stat()
	checkError(e)
	e = sendFileSize(stat.Size(), conn)
	checkError(e)

	fmt.Printf("Sending %v bytes...\n", stat.Size())
	i, e := io.Copy(conn, f)
	checkError(e)
	fmt.Printf("Succesfully sent %v bytes of data\n", i)
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

func archive(name string, p string) (*os.File, error) {
	f, e := os.Create(name)
	if e != nil {
		return f, e
	}
	w := zip.NewWriter(f)
	e = recursiveArchive(p, w)
	e = w.Close()
	return f, e
}

func recursiveArchive(path string, w *zip.Writer) error {
	file, e := os.Open(path)
	if e != nil {
		return e
	}
	defer file.Close()
	s, e := file.Stat()
	if e != nil {
		return e
	}

	if filepath.Ext(file.Name()) == "zip" || !s.IsDir() {
		wr, e := w.Create(path)
		if e != nil {
			return e
		}
		_, e = io.Copy(wr, file)
		return e
	}
	list, e := file.Readdir(0)
	if e != nil {
		return e
	}
	for _, ss := range list {
		p := filepath.Join(path, ss.Name())
		e = recursiveArchive(p, w)
		if e != nil {
			return e
		}
	}
	return nil
}