package target_test

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

type Ping struct {
}

func NewPing() Ping {
	//return new(Ping)
	return Ping{}
}

/*func NewPing() *Ping {
	return new(Ping)
}*/

func (r *Ping) Test() string {
	host := "luutia.com"
	pinger, err := ping.NewPinger(host)

	if err != nil {
		fmt.Println("ERROR:", err)
		return ""
	}

	pinger.Count = 3
	pinger.Run()
	stats := pinger.Statistics()
	fmt.Println("___stats: ", stats)

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%v time=%v ms\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt)
	}
	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = 1
	pinger.Size = 24
	pinger.TTL = 64
	pinger.Timeout = time.Second * 10
	pinger.Interval = time.Second
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		fmt.Println("____ Failed to ping target host:", err)
	}
	return "ok"
}
