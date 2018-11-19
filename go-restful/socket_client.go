package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	var (
		host = "127.0.0.1"
		port = "8080"
		remote = host + ":" + port
		data = make([]byte, 1024)
	)
	var str string
	var msg = make([]byte, 1024)

	conn, err := net.Dial("tcp", remote)
	defer conn.Close()
	if err != nil {
		fmt.Println("server not found")
		os.Exit(-1)
	}

	fmt.Println("Connection OK.")
	for{
		var str string
		fmt.Printf("Enter a sentence:")
		fmt.Scanf("%s", &str)
		if str == "quit" {
			fmt.Println("Communication terminated.")
			os.Exit(1)
		}

		in, err := conn.Write([]byte(str))
		if err != nil {
			fmt.Printf("Error when send to server: %d\n", in)
			os.Exit(0)
		}

		length, err := conn.Read(msg)
		if err != nil {
			fmt.Printf("Error when read from server.\n")
			os.Exit(0)
		}
		str = string(msg[0:length])
		fmt.Println(str)
	}
}