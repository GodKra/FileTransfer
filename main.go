package main

import (
	"flag"
	"fmt"
	"github.com/GodKra/FileTransfer/downloader"
	"github.com/GodKra/FileTransfer/sender"
	"log"
	"net"
	"strings"
)

var (
	// Sender flags: path, ipFlag
	// Downloader flag: name
	path   = flag.String("path", "FileTransfer", "Optional for sender. Usage : -path <path> | eg : -path FileTransfer/test")
	ipFlag = flag.String("ip", "localhost:5555", "Optional for sender. Usage : -ip <ip address:port> | eg : -ip localhost:5555")
	name   = flag.String("fileName", "ftdownload", "Optional for downloader. Usage : -file <FileName> eg : -file test")
	typ    = flag.String("type", "downloader", "Usage: -type <sender/downloader> | eg: -type downloader")
)

func main() {
	flag.Parse()
	switch strings.ToLower(*typ) {
	case "sender":
		conn, e := net.Dial("tcp", *ipFlag)
		checkError(e)
		defer conn.Close()
		sendr := sender.Sender{Connection: conn, Path: *path}
		size, e := sendr.SendFile()
		checkError(e)
		fmt.Printf("Succesfully sent %v bytes of data\n", size)
	case "downloader":
		l, e := net.Listen("tcp", ":5555")
		checkError(e)
		downloadr := downloader.Downloader{Listener: l, Name: *name}
		e = downloadr.DownloadFile()
		checkError(e)
	default:
		log.Fatal("Give a valid 'type'. Either sender or downloader. Sender sends files. Downloader downloads them.")
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
