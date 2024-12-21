package device

import (
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"net"
	"wol-e/internal/logger"
	"wol-e/internal/wol"
)

type Device struct {
	Name string
	Ip   string
	Mac  string
}

func (d Device) GenerateBotText() (string, error) {
	status, err := d.CheckOnline()
	if err != nil {
		return "", err
	}
	text := "name: " + d.Name + "\n" +
		"ip/hostname: " + d.Ip + "\n" +
		"mac: " + d.Mac + "\n"
	if status == true {
		text += "status: 🔋"
	} else {
		text += "status: 🪫"
	}
	return text, nil
}

func (d Device) CheckOnline() (bool, error) {
	pinger, err := probing.NewPinger(d.Ip)
	if err != nil {
		logger.Log.Errorf("pinger create error: %v", err)
		return false, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 3

	err = pinger.Run()
	if err != nil {
		logger.Log.Errorf("ping error: %v", err)
		return false, err
	}
	stats := pinger.Statistics()

	logger.Log.Infof("%s - sent: %d, recieved: %d", d.Ip, stats.PacketsSent, stats.PacketsRecv)

	return stats.PacketsRecv > 0, nil
}

func (d Device) TurnOn() error {
	// The address to broadcast to is usually the default `255.255.255.255` but
	// can be overloaded by specifying an override in the CLI arguments.
	bcastAddr := fmt.Sprintf("%s:%d", d.Ip, 9)
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		logger.Log.Errorf("resolve udp error: %v", err)
		return err
	}

	// Build the magic packet.
	mp, err := wol.New(d.Mac)
	if err != nil {
		logger.Log.Errorf("new wol error: %v", err)
		return err
	}

	// Grab a stream of bytes to send.
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}

	// Grab a UDP connection to send our packet of bytes.
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logger.Log.Errorf("dial udp error: %v", err)
		return err
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			logger.Log.Errorf("close udp error: %v", err)
		}
	}(conn)

	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}

	logger.Log.Infof("Magic packet sent successfully to %s", d.Mac)
	return nil
}
