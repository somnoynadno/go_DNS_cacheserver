package main

import (
	"fmt"
	"go_DNS_cacheserver/parser"
	"net"
)

func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	size, addr, err := conn.ReadFromUDP(buf[:])
	if err != nil {
		return
	}

	fmt.Println(addr.IP, addr.Port)

	dns  := parser.ParseDNS(buf[:size])
	ans  := parser.InflateDNS(dns)
	data := parser.Proxy(ans)

	_, err = conn.WriteToUDP(data, addr)
}

func main() {
	addr := net.UDPAddr{
		Port: 5153,
		IP:   net.ParseIP("0.0.0.0"),
	}

	fmt.Println("Listening on", addr.Port)

	l, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	for {
		handleClient(l)
	}

}