package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepted connected at: ", conn.RemoteAddr())
	buffer := make([]byte, 1024)

	index := 0

	for {
		r, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		message := string(buffer[:r])

		if err := os.MkdirAll("services", 0755); err != nil {
			fmt.Println("folder not found!")
		}

		filename := fmt.Sprintf("service-%d.txt", index)
		dst, err := os.Create(filepath.Join("services", filename))
		if err != nil {
			fmt.Println(err)
		}
		defer dst.Close()
		if _, err := dst.Write(buffer); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Echo: ", message)

		index++
	}
}
