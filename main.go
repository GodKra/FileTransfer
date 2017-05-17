package main

import (
	"flag"
	"fmt"
	"github.com/GodKra/FileTransfer/downloader"
	"github.com/GodKra/FileTransfer/sender"
	"log"
	"net"
	"strings"
	"time"
)

var (
	// Sender flags: path, ipFlag
	// Downloader flag: name
	path   = flag.String("path", "FileTransfer", "Optional for sender. Usage : -path <path> | eg : -path FileTransfer/test")
	ipFlag = flag.String("ip", "0:5151", "Optional for sender. 0:5151 for automatic . Usage : -ip <ip address:port> | eg : -ip localhost:5555")
	name   = flag.String("fileName", "ftdownload", "Optional for downloader. Usage : -file <FileName> eg : -file test")
	typ    = flag.String("type", "downloader", "Usage: -type <sender/downloader> | eg: -type downloader")
)

func main() {
	flag.Parse()
	switch strings.ToLower(*typ) {
	case "sender":
		var conn net.Conn
		switch *ipFlag {
		case "0:5151":
			for i := 0; i < 255; i++ {
				var e error
				ip := fmt.Sprintf("192.168.1.%v:5151", i)
				conn, e = net.DialTimeout("tcp", ip, time.Millisecond*200)
				if e != nil {
					continue
				} else {
					if isDownloaderServer(conn) {
						fmt.Printf("Found IP: %v\n", conn.RemoteAddr())
						break
					}
				}
			}
		default:
			var e error
			conn, e = net.Dial("tcp", *ipFlag)
			checkError(e)
			buf := [1]byte{}
			conn.Read(buf[:])
			if len(buf) > 0 {
				if buf[0] != '`' {
					log.Fatal("Invalid IP address")
				}
			}
		}
		defer conn.Close()

		sendr := sender.Sender{Connection: conn, Path: *path}
		size, e := sendr.SendFile()
		checkError(e)
		fmt.Printf("Succesfully sent %v bytes of data\n", size)
	case "downloader":
		l, e := net.Listen("tcp", ":5151")
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

// isDownloaderServer listens for a response from the server. If it recieves the message "`", then it is a
// downloader's server.
func isDownloaderServer(conn net.Conn) bool {
	buf := [1]byte{}
	conn.Read(buf[:])
	if len(buf) > 0 {
		if buf[0] == '`' {
			return true
		}
	}
	return false
}
