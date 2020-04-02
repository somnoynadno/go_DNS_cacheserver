package DNScache

import (
	"net"
	"time"
)

func ProxyForward(data []byte) ([]byte, error) {
	NS := "8.8.8.8:53"

	s, err := net.ResolveUDPAddr("udp4", NS)
	c, err := net.DialUDP("udp4", nil, s)

	if err != nil {
		panic(err)
	}

	deadline := time.Now().Add(1 * time.Second)
	_ = c.SetReadDeadline(deadline)

	defer c.Close()

	_, err = c.Write(data)

	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	n, _, err := c.ReadFromUDP(buffer)

	if err != nil {
		return nil, err
	}

	return buffer[0:n], nil
}
