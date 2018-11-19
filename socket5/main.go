package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	VER_SOCKET5  = 0x05
	NORMAL       = 0x00
	INDEX_ATYP  = 3
	ATYPE_IPV4   = 0x01
	ATYPE_DOMAIN = 0X03
)

func main() {
	listener, err := net.Listen("tcp", ":1080")
	CheckErr(err)

	// 监听服务，然后对每一个客户端连接都开启一个 coroutine 进一步回应
	for {
		conn, err := listener.Accept()
		CheckErr(err)
		go HandleClientRequest(conn)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}


func HandleClientRequest(client net.Conn) {
	if client == nil {
		return
	}

	// 将客户端发来的第一段数据包读取到一个 1024 字节的 buffer 中（通常协议请求不会超过这个容量）
	var buf [1024]byte
	_, err := client.Read(buf[:])
	CheckErr(err)

	if buf[0] == VER_SOCKET5 {
		client.Write([]byte{VER_SOCKET5, NORMAL})
		n, err := client.Read(buf[:])
		CheckErr(err)
		var host, port string
		switch buf[INDEX_ATYP] {
		case ATYPE_IPV4:
			// 截取 ipv4 地址
			host = net.IPv4(buf[INDEX_ATYP+1], buf[INDEX_ATYP+2],
				buf[INDEX_ATYP+3], buf[INDEX_ATYP+4]).String()
			log.Printf("Client requests ipv4:%s", host)
		case ATYPE_DOMAIN:
			// 截取完整域名
			// 此处也可以利用 buf[4] （表示域名的长度）截取:
			// begin, end := INDEX_ATYP+2, INDEX_ATYP+2+buf[INDEX_ATYP+1]
			begin, end := INDEX_ATYP+2, n-2
			host = string(buf[begin: end])
			log.Printf("Client requests domian name:%s", host)
		}

		port = strconv.Itoa(int(binary.BigEndian.Uint16(buf[n-2:n])))
		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Println(err)
			return
		}

		client.Write([]byte{})
		defer server.Close()

		go io.Copy(client, server)
		io.Copy(server, client)
	}
}
