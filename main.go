package main

import (
	"encoding/binary"
	"fmt"
	"go_DNS_cacheserver/DNScache"
	"net"
)

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	var dns DNScache.DNS

	size, addr, err := conn.ReadFromUDP(buf[:])
	if err != nil {
		panic(err)
	}

	fmt.Println(addr.IP, addr.Port)

	dns = DNScache.ParseDNSRequest(buf[:size])
	answers, err := DNScache.CheckAnswerInCache(dns.Query)

	if err != nil {
		panic(err)
	}

	var data []byte

	if len(answers) == 0 {
		fmt.Println("Using proxy:")

		data = DNScache.ProxyForward(buf[:size])
		dns  = DNScache.ParseDNSResponse(data)

		err = DNScache.PutAnswersInCache(dns.Query, dns.Answers)

		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("From cache:")

		bs := make([]byte, 2)
		var arr [2]byte

		binary.BigEndian.PutUint16(bs, uint16(len(answers)))
		copy(arr[:], bs[:2])

		dns.Header.ARC = arr

		dns.Answers = answers
		data = DNScache.InflateDNS(dns)
	}

	fmt.Println(data)
	_, err = conn.WriteToUDP(data, addr)
}

func main() {
	addr := net.UDPAddr{
		Port: 5153,
		IP:   net.ParseIP("0.0.0.0"),
	}

	err := DNScache.RedisHealthCheck()
	if err != nil {
		panic("Can not connect to redis server")
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