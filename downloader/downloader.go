package downloader

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

type Downloader struct {
	Listener net.Listener
	Conn     net.Conn
	Name     string
}

func (d Downloader) DownloadFile() error {
	e := d.createDir()
	if e != nil {
		return e
	}
	var i int
	for {
		i++
		fmt.Printf("\n-- File %v --\n", i)
		conn, e := d.Listener.Accept()
		if e != nil {
			return e
		}
		d.Conn = conn
		size, e := d.getSize()
		if e != nil {
			return e
		}

		fileName := fmt.Sprintf("%v%v.zip", d.Name, i)
		f, e := os.Create(fileName)
		if e != nil {
			return e
		}
		fmt.Printf("Copying file from %v\n", conn.RemoteAddr())
		e = d.download(f, conn, size)
		if e != nil {
			return e
		}
		fmt.Print("\n")
		f.Close()
	}
}

func (d Downloader) createDir() error {
	os.Mkdir("FileTransfer", 0666)
	e := os.Chdir("FileTransfer")
	if e != nil {
		return e
	}
	return nil
}

func (d Downloader) download(dst io.Writer, src io.Reader, s uint64) error {
	duration := time.Now()
	buf := make([]byte, 32*1024)
	var written int
	const length = 50
	progressbar := bytes.Repeat([]byte{'-'}, length)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			written += nw
			percentage := float64(written) / float64(s) * 100
			for i := 0; i < int(length*float64(percentage/100)); i++ {
				progressbar[i] = '='
			}
			fmt.Printf("\r%v/%v [%s] %.3v%%      ", size(written), size(s), progressbar, percentage)
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

type size int64

func (s size) String() string {
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

func (d Downloader) getSize() (size uint64, e error) {
	b := [8]byte{}
	_, e = d.Conn.Read(b[:])
	size = binary.BigEndian.Uint64(b[:])
	return
}
