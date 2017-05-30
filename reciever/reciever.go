package reciever

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/GodKra/FileTransfer/format"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"time"
)

// Reciever opens a server for the Sender to connect to. Reciever downloads files sent
// by the Sender using RecieverFile() method of Reciever.
type Reciever struct {
	Listener net.Listener
	Conn     net.Conn
	Name     string
}

var ErrInvalidClient = "data recieved from an invalid client"

// RecieveFile creates a Directory using createDir() function then accepts a connection from the listener of
// Reciever. Then it reads the file size using getSize() function. Then reads the rest of the data sent
// by the Sender using download() function
func (r *Reciever) RecieveFile() error {
	e := r.createDir()
	if e != nil {
		return e
	}
	var i int
	for {
		i++
		fmt.Printf("\n-- File %v --\n", i)
		var e error
		r.Conn, e = r.Listener.Accept()
		if e != nil {
			return e
		}
		r.Conn.Write([]byte{'`'})
		if !isSenderClient(r.Conn) {
			r.Conn.Write([]byte{84, 70})
			return fmt.Errorf("%v", ErrInvalidClient)
		}
		r.Conn.Write([]byte{70, 84})
		size, e := r.getSize()
		if e != nil {
			return e
		}

		fileName := fmt.Sprintf("%v%v.zip", r.Name, i)
		f, e := os.Create(fileName)
		if e != nil {
			return e
		}
		fmt.Printf("Copying file from %v\n", r.Conn.RemoteAddr())
		e = r.download(f, size)
		if e != nil {
			return e
		}
		fmt.Print("\n")
		f.Close()
	}
}

// createDir creates a new file called "FileTransfer" (if it exists, does nothing) and Changes working directory
// to "FileTransfer" directory
func (d Reciever) createDir() error {
	os.Mkdir("FileTransfer", 0666)
	e := os.Chdir("FileTransfer")
	if e != nil {
		return e
	}
	return nil
}

// download reads the the data sent by the client (which is data of the file sent). While reading the data,
// it also makes a progress bar and updates it each 32KBs read. When downloading is finished, It prints the
// time taken to download the file.
func (d Reciever) download(file *os.File, s uint64) error {
	duration := time.Now()
	buf := make([]byte, 32*1024)
	var written int
	const length = 50
	progressbar := bytes.Repeat([]byte{'-'}, length)
	for {
		nr, er := d.Conn.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			written += nw
			percentage := float64(written) / float64(s) * 100
			for i := 0; i < int(length*float64(percentage/100)); i++ {
				progressbar[i] = '='
			}
			fmt.Printf("\r%v/%v [%s] %.3v%%      ", format.SizeFormat(written), format.SizeFormat(s), progressbar, percentage)
			if ew != nil {
				return ew
			}
			if nr != nw {
				log.Fatal(io.ErrShortWrite)
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			return er
		}
	}
	timeTaken := time.Now().Sub(duration)
	fmt.Printf("Time taken: %v", formatDuration(timeTaken))
	return nil
}

func formatDuration(d time.Duration) string {
	str := d.String()
	r, _ := regexp.Compile(`([^\d\.]+)`)
	return r.ReplaceAllString(str, "$1 ")
}

// getSize reads the first 8 bytes sent by the Sender (which is the size of the file) and returns it.
func (d Reciever) getSize() (size uint64, e error) {
	b := [8]byte{}
	_, e = d.Conn.Read(b[:])
	size = binary.BigEndian.Uint64(b[:])
	return
}

// isSenderClient is used to validate the Client which is sending data to the Reciever.
func isSenderClient(conn net.Conn) bool {
	buf := [2]byte{}
	conn.Read(buf[:])
	if len(buf) > 1 {
		if buf[0] == 70 && buf[1] == 84 {
			return true
		}
	}
	return false
}
