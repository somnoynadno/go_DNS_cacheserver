package parser

import (
	"fmt"
	"net"
)

type DNS struct {
	ID       [2]byte
	flags    [2]byte
	QC       [2]byte
	ARC      [2]byte
	NSC      [2]byte
	AddRC    [2]byte
	query      DNSQuery
	answers  []DNSAnswer
	info       DNSInfo
}

type DNSQuery struct {
	qType [2]byte
	class [2]byte
	name   []byte
}

type DNSInfo struct {
	additional []byte
}

type DNSAnswer struct {
	name  [2]byte
	qType [2]byte
	class [2]byte
	ttl   [4]byte
	dSize [2]byte
	addr   []byte
}

func ParseDNS(buf []byte) DNS {
	var dns   DNS

	dns.ID    = [2]byte { buf[0], buf[1] }
	dns.flags = [2]byte { buf[2], buf[3] }
	dns.QC    = [2]byte { buf[4], buf[5] }
	dns.ARC   = [2]byte { buf[6], buf[7] }
	dns.NSC   = [2]byte { buf[8], buf[9] }
	dns.AddRC = [2]byte { buf[10], buf[11] }

	dns.query = parseQuery(buf[12:])
	querySize := 4 + len(dns.query.name)

	dns.info.additional = buf[12+querySize:]

	return dns
}

func parseDNSResponse(buf []byte) DNS {
	var dns DNS

	dns.ID    = [2]byte { buf[0],  buf[1]  }
	dns.flags = [2]byte { buf[2],  buf[3]  }
	dns.QC    = [2]byte { buf[4],  buf[5]  }
	dns.ARC   = [2]byte { buf[6],  buf[7]  }
	dns.NSC   = [2]byte { buf[8],  buf[9]  }
	dns.AddRC = [2]byte { buf[10], buf[11] }

	dns.query = parseQuery(buf[12:])
	querySize := 4 + len(dns.query.name)

	padding := 12 + querySize

	for i := 0; i < int(dns.ARC[1]); i++ {
		answer := parseDNSAnswer(buf[padding:])
		dns.answers = append(dns.answers, answer)
		padding += 12 + len(answer.addr)
		fmt.Println("Answer:", answer)
	}

	dns.info.additional = buf[padding:]

	return dns
}

func parseDNSAnswer(data []byte) DNSAnswer {
	var answer DNSAnswer

	answer.name  = [2]byte { data[0],  data[1]  }
	answer.qType = [2]byte { data[2],  data[3]  }
	answer.class = [2]byte { data[4],  data[5]  }
	answer.ttl   = [4]byte { data[6],  data[7], data[8], data[9] }
	answer.dSize = [2]byte { data[10], data[11] }

	s := int(answer.dSize[1])
	answer.addr  = data[12:12+s]

	return answer
}

func parseQuery(data []byte) DNSQuery {
	var query DNSQuery

	for i, v := range data {
		if v == 0 {
			query.name = data[:i+1]
			query.qType = [2]byte { data[i+1], data[i+2] }
			query.class = [2]byte { data[i+3], data[i+4] }
			break
		}
	}

	return query
}

func InflateDNS(dns DNS) []byte {
	var bytes []byte

	bytes = append(bytes, dns.ID[:]...)
	bytes = append(bytes, dns.flags[:]...)
	bytes = append(bytes, dns.QC[:]...)
	bytes = append(bytes, dns.ARC[:]...)
	bytes = append(bytes, dns.NSC[:]...)
	bytes = append(bytes, dns.AddRC[:]...)
	bytes = append(bytes, dns.query.name[:]...)
	bytes = append(bytes, dns.query.qType[:]...)
	bytes = append(bytes, dns.query.class[:]...)

	for _, a := range dns.answers {
		bytes = append(bytes, a.name[:]...)
		bytes = append(bytes, a.qType[:]...)
		bytes = append(bytes, a.class[:]...)
		bytes = append(bytes, a.ttl[:]...)
		bytes = append(bytes, a.dSize[:]...)
		bytes = append(bytes, a.addr[:]...)
	}

	bytes = append(bytes, dns.info.additional[:]...)

	fmt.Println(bytes)
	return bytes
}

func Proxy(data []byte) []byte {
	NS := "8.8.8.8:53"

	s, err := net.ResolveUDPAddr("udp4", NS)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The proxy server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	_, err = c.Write(data)

	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)
	n, _, err := c.ReadFromUDP(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Reply:", buffer[0:n])
	dns := parseDNSResponse(buffer[0:n])
	fmt.Println("Parsed:", dns)

	return buffer
}
