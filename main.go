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
	// Reciever flag: name
	// Both: port
	filePath = flag.String("filePath", "", "Usage : -filePath <path> | eg : -filePath FileTransfer/test")
	ipFlag   = flag.String("ip", "0", "0 for automatic . Usage : -ip <ip address> | eg : -ip localhost")
	port     = flag.String("port", "7084", "Usage: -port <port> | eg: --port 1234")
	name     = flag.String("saveName", "ftdownload", "Usage : -fileName <FileName> eg : -fileName test")
	typ      = flag.String("type", "reciever", "Optional Usage: -type <sender/reciever> | eg: -type reciever")
	help     = flag.Bool("help", false, "Show help. Usage: --help")
)

func main() {
	flag.Parse()
	if *help {
		printHelp()
		return
	}
	if *filePath != "" {
		*typ = "sender"
	}
	switch strings.ToLower(*typ) {
	case "sender":
		var conn net.Conn
		switch *ipFlag {
		case "0":
			for i := 0; i < 255; i++ {
				var e error
				ip := fmt.Sprintf("192.168.1.%v:%v", i, *port)
				conn, e = net.DialTimeout("tcp", ip, time.Millisecond*200)
				if e != nil {
					continue
				} else {
					if isRecieverServer(conn) {
						fmt.Printf("Found IP: %v\n", conn.RemoteAddr())
						// Sent to validate Sender server from reciever server
						conn.Write([]byte{70, 84})
						buf := [2]byte{}
						conn.Read(buf[:])
						if buf[0] == 84 && buf[1] == 70 {
							log.Fatal("Reciever identifies you as an invalid client")
						}
						break
					}
				}
			}
		default:
			var e error
			conn, e = net.Dial("tcp", fmt.Sprintf("%v:%v", *ipFlag, *port))
			checkError(e)
			if !isRecieverServer(conn) {
				log.Fatal("Invalid IP address")
			}
			// Sent to validate Sender server from reciever server
			conn.Write([]byte{70, 84})
		}
		defer conn.Close()

		sendr := sender.Sender{Connection: conn, FilePath: *filePath}
		size, e := sendr.SendFile()
		checkError(e)
		fmt.Printf("Succesfully sent %v bytes of data\n", size)
	case "reciever":
		l, e := net.Listen("tcp", fmt.Sprintf(":%v", *port))
		checkError(e)
		recievr := reciever.Reciever{Listener: l, Name: *name}
		e = recievr.RecieveFile()
		if e.Error() == reciever.ErrInvalidClient {
			log.Fatalf("Invalid Client Sent Data: %v", recievr.Conn.RemoteAddr())
		}
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
	fmt.Println("\t--filePath [value]: Path of the file you want to transfer. Must for Sender. If this exists, type will be considered as a sender")
	fmt.Println("\t--ip [value]:       The IP of the Reciever you want to send the file to. 0 for automatic. Optional for Sender")
	fmt.Println("\t--saveName [value]: Name to use when saving the recieved files. Optional for Reciever")
	fmt.Println("\t--type [value]:     Optional. The type of filetransfer. 'sender' to send files. 'reciever' to recieve files")
	fmt.Println("\t--help:             Prints this.")
	fmt.Println("\t--port [value]:     The port of the Reciever. Default is '7084'. Optional for both Sender and Reciever")
}
