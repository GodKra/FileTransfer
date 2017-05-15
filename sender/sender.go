package sender

import (
	"archive/zip"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

const zipName = "temp.zip"

type Sender struct {
	Connection net.Conn
	Path       string
}

func (s Sender) SendFile() (size int64, e error) {
	fmt.Println("Zipping File..")
	f, e := archive(s.Path, zipName)
	if e != nil {
		return 0, e
	}
	defer os.Remove(zipName)
	defer f.Close()
	_, e = f.Seek(0, io.SeekStart)
	if e != nil {
		return 0, e
	}
	stat, e := f.Stat()
	if e != nil {
		return 0, e
	}
	e = sendFileSize(stat.Size(), s.Connection)
	if e != nil {
		return 0, e
	}

	fmt.Printf("Sending %v bytes...\n", stat.Size())
	i, e := io.Copy(s.Connection, f)
	if e != nil {
		return 0, e
	}
	return i, nil
}

func sendFileSize(size int64, conn net.Conn) error {
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], uint64(size))
	_, e := conn.Write(b[:])
	return e
}

func archive(path string, name string) (*os.File, error) {
	f, e := os.Create(name)
	if e != nil {
		return f, e
	}
	w := zip.NewWriter(f)
	e = recursiveArchive(path, w)
	e = w.Close()
	return f, e
}

func recursiveArchive(path string, w *zip.Writer) error {
	file, e := os.Open(path)
	if e != nil {
		return e
	}
	defer file.Close()
	s, e := file.Stat()
	if e != nil {
		return e
	}

	if filepath.Ext(file.Name()) == "zip" || !s.IsDir() {
		wr, e := w.Create(path)
		if e != nil {
			return e
		}
		_, e = io.Copy(wr, file)
		return e
	}
	list, e := file.Readdirnames(0)
	if e != nil {
		return e
	}
	for _, ss := range list {
		p := filepath.Join(path, ss)
		e = recursiveArchive(p, w)
		if e != nil {
			return e
		}
	}
	return nil
}
