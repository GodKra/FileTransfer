package main

import (
	"net"
	"log"
	"fmt"
    "os"
    "io"
    "encoding/binary"
)

func main() {
    args := os.Args
    args = args[1:]
    file := args[0]
    ip := args[1]
    f, e := os.Open(file)
    checkError(e)

    conn := connect(ip)
    stat, e := f.Stat()
    checkError(e)
    conn.Write([]byte(stat.Name() + "\n"))
    sendFileSize(stat.Size(), conn)

    i, e := io.Copy(conn, f)
    defer checkError(e)
    fmt.Println("Sending...")
    fmt.Printf("succesfully sent %v bytes of data\n", i)
    conn.Close()
}

func connect(ip string) net.Conn{
    conn, e := net.Dial("tcp", ip)
    checkError(e)
    return conn
}

func checkError(e error) {
    if e != nil {
        log.Fatal(e)
    }
} 

func sendFileSize(size int64, conn net.Conn) {
    b := [8]byte{}
    binary.BigEndian.PutUint64(b[:], uint64(size))
    conn.Write(b[:])
}