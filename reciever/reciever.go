package reciever

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// Reciever opens a server for the Sender to connect to. Reciever downloads files sent
// by the Sender using RecieverFile() method of Reciever.
type Reciever struct {
	Listener net.Listener
	Conn     net.Conn
	Name     string
}

// RecieveFile creates a Directory using createDir() function then accepts a connection from the listener of
// Reciever. Then it reads the file size using getSize() function. Then reads the rest of the data sent
// by the Sender using download() function
func (d *Reciever) RecieveFile() error {
	e := d.createDir()
	if e != nil {
		return e
	}
	var i int
	for {
		i++
		fmt.Printf("\n-- File %v --\n", i)
		var e error
		d.Conn, e = d.Listener.Accept()
		if e != nil {
			return e
		}
		d.Conn.Write([]byte{'`'})
		size, e := d.getSize()
		if e != nil {
			return e
		}

		fileName := fmt.Sprintf("%v%v.zip", d.Name, i)
		f, e := os.Create(fileName)
		if e != nil {
			return e
		}
		fmt.Printf("Copying file from %v\n", d.Conn.RemoteAddr())
		e = d.download(f, size)
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
			fmt.Printf("\r%v/%v [%s] %.3v%%      ", sizeFormat(written), sizeFormat(s), progressbar, percentage)
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
	fmt.Printf("Time taken: %v", time.Now().Sub(duration))
	return nil
}

// This type is used to change the way how fmt.Print prints the file size.
type sizeFormat int64

func (s sizeFormat) String() string {
	switch {
	case s < 1<<10:
		return fmt.Sprintf("%v B", float64(s))
	case s < 1<<20:
		return fmt.Sprintf("%.2f KB", float64(s)/(1<<10))
	case s < 1<<30:
		return fmt.Sprintf("%.2f MB", float64(s)/(1<<20))
	default:
		return fmt.Sprintf("%.2f GB", float64(s)/(1<<30))
	}
}

// getSize reads the first 8 bytes sent by the Sender (which is the size of the file) and returns it.
func (d Reciever) getSize() (size uint64, e error) {
	b := [8]byte{}
	_, e = d.Conn.Read(b[:])
	size = binary.BigEndian.Uint64(b[:])
	return
}
