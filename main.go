package main

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"go_DNS_cacheserver/DNScache"
	"net"
)

var log = logrus.New()

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	var dns DNScache.DNS

	size, addr, err := conn.ReadFromUDP(buf[:])
	if err != nil {
		log.Error(err)
		return
	}

	log.WithFields(logrus.Fields{
		"IP": addr.IP,
		"port": addr.Port,
	}).Info("UDP connection")

	dns = DNScache.ParseDNSRequest(buf[:size])
	answers, err := DNScache.CheckAnswerInCache(dns.Query)

	if err != nil {
		log.Error(err)
		return
	}

	var data []byte

	if len(answers) == 0 {
		log.Info("Using proxy")

		data, err = DNScache.ProxyForward(buf[:size])
		if err != nil {
			log.Error(err)
			return
		}
		dns  = DNScache.ParseDNSResponse(data)

		err = DNScache.PutAnswersInCache(dns.Query, dns.Answers)
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		log.Info("From cache")

		bs := make([]byte, 2)
		var arr [2]byte

		binary.BigEndian.PutUint16(bs, uint16(len(answers)))
		copy(arr[:], bs[:2])

		dns.Header.ARC = arr
		dns.Answers = answers

		data = DNScache.InflateDNS(dns)
	}

	log.WithFields(logrus.Fields{
		"data": data,
	}).Debug("Success")

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

	l, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.Info("Listening on: ", addr.Port)
	for {
		handleClient(l)
	}

}