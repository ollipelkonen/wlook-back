package target_test

import (
	"fmt"

	"github.com/go-ping/ping"
	tgtg "github.com/ollipelkonen/wlook_back/target_test/defs"
)

type Ping struct {
}

func (r *tgtg.Target_test) Test() string {
	host := "luutia.com"
	pinger, err := ping.NewPinger(host)

	if err != nil {
		fmt.Println("ERROR:", err)
		return ""
	}
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.Ttl)
	}
	pinger.Count = 1
	err = pinger.Run()
	if err != nil {
		fmt.Println("Failed to ping target host:", err)
	}
	return "ok"
}
