package device

import (
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"net"
	"wol-e/internal/logger"
)

type Device struct {
	Name string
	Ip   string
	Mac  string
}

func (d Device) CheckOnline() (bool, error) {
	pinger, err := probing.NewPinger(d.Ip)
	if err != nil {
		return false, err
	}
	pinger.Count = 3

	err = pinger.Run()
	if err != nil {
		return false, err
	}
	stats := pinger.Statistics()

	logger.Log.Infof("%s - sent: %d, recieved: %d", d.Ip, stats.PacketsSent, stats.PacketsRecv)

	return true, nil
}

func (d Device) TurnOn() error {
	// Convert MAC address to byte slice
	mac, err := parseMACAddress(d.Mac)
	if err != nil {
		return err
	}

	// Create a magic packet: 6 x 0xFF + 16 x MAC address
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 6; i < 102; i++ {
		packet[i] = mac[(i-6)%6]
	}

	// Resolve broadcast address
	addr := fmt.Sprintf("%s:9", d.Ip) // Default WOL port is 9

	// Create UDP connection
	conn, err := net.Dial("udp", addr)
	if err != nil {
		logger.Log.Errorf("could not dial UDP connection: %v", err)
		return err
	}
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			logger.Log.Errorf("could not close connection: %v", err)
		}
	}(conn)

	// Send the packet
	_, err = conn.Write(packet)
	if err != nil {
		logger.Log.Errorf("could not send packet: %v", err)
		return err
	}

	logger.Log.Infof("WOL packet sent to %s, %s", d.Ip, addr)
	return nil
}

// parseMACAddress parses the MAC address string (e.g., "00:11:22:33:44:55") to a byte slice
func parseMACAddress(mac string) ([]byte, error) {
	var result []byte
	_, err := fmt.Sscanf(mac, "%02x:%02x:%02x:%02x:%02x:%02x",
		&result[0], &result[1], &result[2], &result[3], &result[4], &result[5])
	if err != nil {
		logger.Log.Errorf("invalid MAC address: %v", err)
		return nil, err
	}
	return result, nil
}
