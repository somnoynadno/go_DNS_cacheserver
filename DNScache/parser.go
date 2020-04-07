package DNScache

import (
	"encoding/binary"
)

type DNS struct {
	Header  DNSHeader
	Query   DNSQuery
	Answers []DNSAnswer
	Info    DNSInfo
}

type DNSHeader struct {
	ID    [2]byte
	Flags [2]byte
	QC    [2]byte
	ARC   [2]byte
	NSC   [2]byte
	AddRC [2]byte
}

type DNSQuery struct {
	QType [2]byte
	Class [2]byte
	Name   []byte
}

type DNSInfo struct {
	Additional []byte
}

type DNSAnswer struct {
	Name   []byte
	QType [2]byte
	Class [2]byte
	TTL   [4]byte
	DSize [2]byte
	Addr   []byte
}

func NewRequest(buf []byte) *DNS {
	dns := new(DNS)

	dns.Header = *newDNSHeader(buf[:12])
	dns.Query  = *newDNSQuery(buf[12:])

	querySize := 4 + len(dns.Query.Name)
	dns.Info.Additional = buf[12+querySize:]

	return dns
}

func NewResponse(buf []byte) *DNS {
	dns := new(DNS)

	dns.Header = *newDNSHeader(buf[:12])
	dns.Query  = *newDNSQuery(buf[12:])

	querySize := 4 + len(dns.Query.Name)
	padding := 12 + querySize

	for i := 0; i < int(binary.BigEndian.Uint16(dns.Header.ARC[:])); i++ {
		answer := *newDNSAnswer(buf[padding:])
		dns.Answers = append(dns.Answers, answer)
		padding += 10 + len(answer.Addr) + len(answer.Name)
	}

	dns.Info.Additional = buf[padding:]

	return dns
}

func newDNSAnswer(data []byte) *DNSAnswer {
	answer := new(DNSAnswer)

	for i, v := range data {
		if v == 192 {
			answer.Name = data[:i+2]

			answer.QType = [2]byte { data[i+2],  data[i+3]  }
			answer.Class = [2]byte { data[i+4],  data[i+5]  }
			answer.TTL   = [4]byte { data[i+6],  data[i+7], data[i+8], data[i+9] }
			answer.DSize = [2]byte { data[i+10], data[i+11] }

			s := binary.BigEndian.Uint16(answer.DSize[:])
			answer.Addr = data[i+12:i+12+int(s)]
			break
		}
	}

	return answer
}

func newDNSQuery(data []byte) *DNSQuery {
	query := new(DNSQuery)

	for i, v := range data {
		if v == 0 {
			query.Name = data[:i+1]
			query.QType = [2]byte {data[i+1], data[i+2] }
			query.Class = [2]byte {data[i+3], data[i+4] }
			break
		}
	}

	return query
}

func newDNSHeader(data []byte) *DNSHeader {
	header := new(DNSHeader)

	header.ID    = [2]byte { data[0],  data[1]  }
	header.Flags = [2]byte { data[2],  data[3]  }
	header.QC    = [2]byte { data[4],  data[5]  }
	header.ARC   = [2]byte { data[6],  data[7]  }
	header.NSC   = [2]byte { data[8],  data[9]  }
	header.AddRC = [2]byte { data[10], data[11] }

	return header
}
