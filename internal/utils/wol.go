package utils

import (
	"fmt"
	"net"
)

// 发送wol包
func sendWOLPacket(nasMac string) error {

	mac, err := net.ParseMAC(nasMac)
	if err != nil {
		return err
	}

	// Create the magic packet
	var packet = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	for i := 0; i < 16; i++ {
		packet = append(packet, mac...)
	}

	// Send the magic packet to the broadcast address
	conn, err := net.Dial("udp", net.IPv4bcast.String()+":6")

	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(packet)

	if err != nil {
		return err
	}

	return nil
}

// wol
func WakeOnLanHandler(mac string) bool {

	err := sendWOLPacket(mac)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
