package main

import (
	"flag"
	"fmt"
	"github.com/GodKra/FileTransfer/reciever"
	"github.com/GodKra/FileTransfer/sender"
	"log"
	"net"
	"strings"
	"time"
)

var (
	// Sender flags: filePath, ipFlag
	// Downloader flag: name
	filePath = flag.String("filePath", "", "Usage : -filePath <path> | eg : -filePath FileTransfer/test")
	ipFlag   = flag.String("ip", "0", "0 for automatic . Usage : -ip <ip address:port> | eg : -ip localhost:5555")
	name     = flag.String("fileName", "ftdownload", "Usage : -fileName <FileName> eg : -fileName test")
	typ      = flag.String("type", "reciever", "Usage: -type <sender/reciever> | eg: -type downloader")
	help     = flag.Bool("help", false, "Show help. Usage: --help")
)

func main() {
	flag.Parse()
	if *help {
		printHelp()
		return
	}
	switch strings.ToLower(*typ) {
	case "sender":
		var conn net.Conn
		switch *ipFlag {
		case "0":
			for i := 0; i < 255; i++ {
				var e error
				ip := fmt.Sprintf("192.168.1.%v:5151", i)
				conn, e = net.DialTimeout("tcp", ip, time.Millisecond*200)
				if e != nil {
					continue
				} else {
					if isRecieverServer(conn) {
						fmt.Printf("Found IP: %v\n", conn.RemoteAddr())
						break
					}
				}
			}
		default:
			var e error
			conn, e = net.Dial("tcp", fmt.Sprintf("%v:5151", *ipFlag))
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

		sendr := sender.Sender{Connection: conn, Path: *filePath}
		size, e := sendr.SendFile()
		checkError(e)
		fmt.Printf("Succesfully sent %v bytes of data\n", size)
	case "reciever":
		l, e := net.Listen("tcp", ":5151")
		checkError(e)
		recievr := reciever.Reciever{Listener: l, Name: *name}
		e = recievr.RecieveFile()
		checkError(e)
	default:
		log.Fatal("Give a valid 'type'. Either sender or reciever. Sender sends files. Reciever recieves them.")
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// isRecieverServer listens for a response from the server. If it recieves the message "`", then it is a
// reciever's server.
func isRecieverServer(conn net.Conn) bool {
	buf := [1]byte{}
	conn.Read(buf[:])
	if len(buf) > 0 {
		if buf[0] == '`' {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Println("Usage: filetransfer <flags>")
	fmt.Println("\nAvailable Flags: ")
	fmt.Println("\t--filePath [value]: Path of the file you want to transfer. Must for Sender")
	fmt.Println("\t--ip [value]:       The IP of the downloader you want to send the file to. 0 for automatic. Optional for Sender")
	fmt.Println("\t--fileName [value]: Name to use when saving the recieved files. Optional for Downloader")
	fmt.Println("\t--type [value]:     The type of filetransfer. 'sender' to send files. 'reciever' to recieve files")
	fmt.Println("\t--help:             Prints this.")
}