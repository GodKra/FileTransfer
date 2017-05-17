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

// Sender sends files to the Downloader using the SendFile() method of Sender.
type Sender struct {
	Connection net.Conn
	Path       string
}

// SendFile first archives a temporary file. Then sends the temp file's size using sendFileSize method. And at last,
// it sends all the data of the file to the Downloader.
func (s Sender) SendFile() (size int64, e error) {
	fmt.Println("Zipping File..")
	f, e := archive(s.Path)
	if e != nil {
		return 0, e
	}
	defer os.Remove(f.Name())
	defer f.Close()
	_, e = f.Seek(0, io.SeekStart)
	if e != nil {
		return 0, e
	}
	stat, e := f.Stat()
	if e != nil {
		return 0, e
	}
	e = s.sendFileSize(stat.Size())
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

// sendFileSize Sends the file size to the Downloader as a Big Endian.
func (s Sender) sendFileSize(size int64) error {
	b := [8]byte{}
	binary.BigEndian.PutUint64(b[:], uint64(size))
	_, e := s.Connection.Write(b[:])
	return e
}

// archive Creates a zip file and archives a directory to it using recursiveArchive function.
func archive(path string) (*os.File, error) {
	f, e := os.Create(zipName)
	if e != nil {
		return f, e
	}
	w := zip.NewWriter(f)
	e = recursiveArchive(path, w)
	e = w.Close()
	return f, e
}

// recursiveArchive opens a file and checks if it is a Directory of a file. If it is a file, it will directly
// archives the data of the file to the zip file. If it is a directory, it will read all the file names in the
// directory and opens them using recursiveArchive.
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
