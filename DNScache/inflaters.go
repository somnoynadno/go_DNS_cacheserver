package DNScache


func (dns DNS) Inflate() []byte {
	var bytes []byte

	bytes = append(bytes, dns.Header.Inflate()...)
	bytes = append(bytes, dns.Query.Inflate()...)

	for _, a := range dns.Answers {
		bytes = append(bytes, a.Inflate()...)
	}

	bytes = append(bytes, dns.Info.Additional[:]...)

	return bytes
}

func (query DNSQuery) Inflate() []byte {
	var bytes []byte

	bytes = append(bytes, query.Name[:]...)
	bytes = append(bytes, query.QType[:]...)
	bytes = append(bytes, query.Class[:]...)

	return bytes
}

func (answer DNSAnswer) Inflate() []byte {
	var bytes []byte

	bytes = append(bytes, answer.Name[:]...)
	bytes = append(bytes, answer.QType[:]...)
	bytes = append(bytes, answer.Class[:]...)
	bytes = append(bytes, answer.TTL[:]...)
	bytes = append(bytes, answer.DSize[:]...)
	bytes = append(bytes, answer.Addr[:]...)

	return bytes
}

func (header DNSHeader) Inflate() []byte {
	var bytes []byte

	bytes = append(bytes, header.ID[:]...)
	bytes = append(bytes, header.Flags[:]...)
	bytes = append(bytes, header.QC[:]...)
	bytes = append(bytes, header.ARC[:]...)
	bytes = append(bytes, header.NSC[:]...)
	bytes = append(bytes, header.AddRC[:]...)

	return bytes
}


