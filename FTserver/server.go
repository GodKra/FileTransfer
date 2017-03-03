package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	
)

func main() {
	name := os.Args[1]
	l, e := listenServer(":5555")
	checkError(e)
	conn, e := l.Accept()
	checkError(e)
	bufcon := bufio.NewReader(conn)
	checkError(e)
	size, e := getSize(bufcon)
	checkError(e)
	os.Mkdir("FileTransfer", 0666)
	os.Chdir("FileTransfer")
	f, _ := os.Create(name + ".zip")
	fmt.Printf("Copying file(%v) from %v\n", name, conn.RemoteAddr())
	copyBuffer(f, conn, size)
	checkError(e)
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getSize(conn io.Reader) (size uint64, e error) {
	b := [8]byte{}
	_, e = conn.Read(b[:])
	size = binary.BigEndian.Uint64(b[:])
	return
}

func listenServer(addr string) (listener net.Listener, e error) {
	listener, e = net.Listen("tcp", addr)
	checkError(e)
	return
}

func copyBuffer(dst io.Writer, src io.Reader, size uint64) {
	buf := make([]byte, 32*1024)
	var written int
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			written += nw
			fmt.Printf("\r%v of %v bytes", written, size)
			checkError(ew)
			if nr != nw {
				log.Fatal(io.ErrShortWrite)
			}
		}
		if er == io.EOF {
			break
		}
		checkError(er)
	}
}