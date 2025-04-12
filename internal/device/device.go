package device

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	probing "github.com/prometheus-community/pro-bing"
	"net"
	"time"
	"wol-e/internal/logger"
	"wol-e/internal/wol"
)

type Device struct {
	Name string
	Ip   string
	Mac  string
	Cron *[]string
}

func (d Device) GenerateBotText() string {
	status, _ := d.CheckOnline()
	text := "name: " + d.Name + "\n" +
		"ip/hostname: " + d.Ip + "\n" +
		"mac: " + d.Mac + "\n"
	if status == true {
		text += "status: 🔋"
	} else if status == false {
		text += "status: 🪫"
	}
	return text
}

func (d Device) CheckOnline() (bool, error) {
	pinger, err := probing.NewPinger(d.Ip)
	if err != nil {
		logger.Log.Errorf("pinger create error: %v", err)
		return false, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 5
	pinger.Timeout = time.Millisecond * 500

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

	logger.Log.Infof("Magic packet sent successfully to %s - %s", bcastAddr, d.Mac)
	return nil
}

func (d Device) InitCron(s gocron.Scheduler) error {
	if d.Cron == nil {
		return nil
	}

	for _, cron := range *d.Cron {
		_, err := s.NewJob(
			gocron.CronJob(cron, false),
			gocron.NewTask(func(d Device) {
				err := d.TurnOn()
				if err != nil {
					logger.Log.Errorf("cron turn on error: %v", err)
					return
				}
			}, d),
		)
		if err != nil {
			return err
		}
		logger.Log.Infof("cron started: %s - \"%s\"", d.Name, cron)
	}

	return nil
}
