package DNScache

import (
	"fmt"
)

func InflateDNS(dns DNS) []byte {
	var bytes []byte

	bytes = append(bytes, InflateDNSHeader(dns.Header)...)
	bytes = append(bytes, InflateDNSQuery(dns.Query)...)

	for _, a := range dns.Answers {
		bytes = append(bytes, InflateDNSAnswer(a)...)
	}

	bytes = append(bytes, dns.Info.Additional[:]...)

	fmt.Println(bytes)
	return bytes
}

func InflateDNSQuery(query DNSQuery) []byte {
	var bytes []byte

	bytes = append(bytes, query.Name[:]...)
	bytes = append(bytes, query.QType[:]...)
	bytes = append(bytes, query.Class[:]...)

	return bytes
}

func InflateDNSAnswer(answer DNSAnswer) []byte {
	var bytes []byte

	bytes = append(bytes, answer.Name[:]...)
	bytes = append(bytes, answer.QType[:]...)
	bytes = append(bytes, answer.Class[:]...)
	bytes = append(bytes, answer.TTL[:]...)
	bytes = append(bytes, answer.DSize[:]...)
	bytes = append(bytes, answer.Addr[:]...)

	return bytes
}

func InflateDNSHeader(header DNSHeader) []byte {
	var bytes []byte

	bytes = append(bytes, header.ID[:]...)
	bytes = append(bytes, header.Flags[:]...)
	bytes = append(bytes, header.QC[:]...)
	bytes = append(bytes, header.ARC[:]...)
	bytes = append(bytes, header.NSC[:]...)
	bytes = append(bytes, header.AddRC[:]...)

	return bytes
}


