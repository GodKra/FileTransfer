package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"bytes"
)

func main() {
	name := flag.String("file", "ftdownload", "Usage : -file <FileName> eg : -file test")
	flag.Parse()
	l, e := listenServer(":5555")
	var i int
	os.Mkdir("FileTransfer", 0666)
	os.Chdir("FileTransfer")
	for {
		i++
		fmt.Printf("\n---- File %v ----\n", i)
		checkError(e)
		conn, e := l.Accept()
		checkError(e)

		bufcon := bufio.NewReader(conn)
		checkError(e)
		size, e := getSize(bufcon)
		checkError(e)

		f, _ := os.Create(*name + ".zip")
		fmt.Printf("Copying file(%v) from %v\n", *name, conn.RemoteAddr())
		download(conn, f, size)
		fmt.Print("\n")
		f.Close()
	}
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

func download(src io.Reader, dst io.Writer, size uint64) {
	buf := make([]byte, 32*1024)
	var written int
	const length = 50
	progressbar := bytes.Repeat([]byte{'-'}, length)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			written += nw
			percentage := float64(written) / float64(size) * 100
			for i := 0; i < int(length * float64(percentage / 100)); i++ {
				progressbar[i] = '='
			}
			fmt.Printf("\r[%s] %v/%v %.3v%%  ", progressbar, written, size, percentage)

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
