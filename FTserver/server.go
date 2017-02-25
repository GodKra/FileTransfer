package main

import (
	"net"
	"log"
	"os"
	"io"
	"fmt"
	"encoding/binary"
	"bufio"
	"strings"
)

func main() {
    l, e := listenServer("localhost:5555")
    checkError(e)
    conn, e := l.Accept()
	checkError(e)
	name := getName(conn)
	size := setSize(conn)
	os.Mkdir("FileTransfer", 0666)
	os.Chdir("FileTransfer")
    f, _ := os.Create(name)
	fmt.Printf("Copying file(%v) from %v\n", name, conn.RemoteAddr())
    copyBuffer(f, conn, size)
    checkError(e)
}

func checkError(e error) {
	if e != nil {
        log.Fatal(e)
    }
}

func setSize(conn net.Conn) uint64{
	b := [8]byte{}
	conn.Read(b[:])
	size := binary.BigEndian.Uint64(b[:])
	return size
}

func getName(conn net.Conn) (name string) {
	r := bufio.NewReader(conn)
	b, _ := r.ReadBytes('\n')
	name = strings.TrimSpace(string(b[:]))
	return name
}

func listenServer(addr string) (listener net.Listener, e error){
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
            fmt.Printf("\r%v of %v", written, size)
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