package DNScache

import (
	"net"
)

func ProxyForward(data []byte) []byte {
	NS := "8.8.8.8:53"

	s, err := net.ResolveUDPAddr("udp4", NS)
	c, err := net.DialUDP("udp4", nil, s)

	if err != nil {
		panic(err)
	}

	defer c.Close()

	_, err = c.Write(data)

	if err != nil {
		return nil
	}

	buffer := make([]byte, 1024)
	n, _, err := c.ReadFromUDP(buffer)

	if err != nil {
		return nil
	}

	return buffer[0:n]
}
