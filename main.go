package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/vxcute/goarping/protocols/arp"
	"github.com/vxcute/goarping/protocols/ethernet"
)

var ( 
	netifF = flag.String("i", "", "network interface")
	targetIPF = flag.String("ip", "", "target ip")
	timeout = flag.Duration("t", time.Second*2, "timeout")
)

const ( 
	BroadcastAddress = "FF:FF:FF:FF:FF:FF"
)

func htons(n uint16) uint16 {
	return (n << 8) & 0xFF00 | (n >> 8) & 0x00FF
}

func GetHostIP() (net.IP, error) {
	
	hostname, err := os.Hostname()
	
	if err != nil {
		return nil, err
	}

	ips, err := net.LookupIP(hostname)

	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if net.ParseIP(ip.String()).To4() != nil {
			return ip, nil
		}
	}

	return nil, errors.New("failed to find ip")
}

func main() {

	flag.Parse()

	if *netifF == "" || *targetIPF == "" {
		flag.Usage()
		os.Exit(1)
	}

	if net.ParseIP(*targetIPF).To4() == nil {
		fmt.Println("[!] IPv6 Is Unsupported yet")
		os.Exit(1)
	}

	netif, err := net.InterfaceByName(*netifF)

	if err != nil {
		log.Fatal(err)
	}

	sockfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ARP)))

	if err != nil {
		log.Fatal(err)
	}

	defer syscall.Close(sockfd)

	destMac, err := net.ParseMAC(BroadcastAddress)

	if err != nil {
		log.Fatal(err)
	}

	hostIP, err := GetHostIP()

	if err != nil {
		log.Fatal(err)
	}

	arpIPv4 := arp.NewARPIPv4(netif.HardwareAddr, destMac, hostIP, net.ParseIP(*targetIPF))
	arpHdr := arp.NewPacket(arp.ArpEthernet, syscall.ETH_P_IP, arp.HardwareSize, arp.ProtocolSize, arp.ArpRequest, arpIPv4.Bytes())
	ethHdr := ethernet.NewFrame(destMac, netif.HardwareAddr, syscall.ETH_P_ARP, arpHdr.Bytes())

	sockAddr := &syscall.SockaddrLinklayer{
		Ifindex: netif.Index,
		Halen:   6,
		Pkttype: syscall.PACKET_BROADCAST,
		Hatype: syscall.ARPHRD_ETHER,
		Protocol: syscall.ETH_P_ARP,
		Addr: [8]byte{
			netif.HardwareAddr[0], netif.HardwareAddr[1], netif.HardwareAddr[2],
			netif.HardwareAddr[3], netif.HardwareAddr[4], netif.HardwareAddr[5], 
			0, 0,
		},
	}

	bindAddr := &syscall.SockaddrLinklayer{
		Ifindex: netif.Index,
	}

	if err := syscall.Bind(sockfd, bindAddr); err != nil {
		log.Fatal(err)
	}

	err = syscall.Sendto(sockfd, ethHdr.Bytes(), 0, sockAddr)

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 100)

	go func() {
		timer := time.NewTimer(*timeout)
		<-timer.C
		fmt.Println("Timeout elapsed")
		os.Exit(1)
	}()

	n, _, err := syscall.Recvfrom(sockfd, buf, 0)

	if err != nil {
		log.Fatal(err)
	}

	senderMAC := arp.ARPIPv4FromBytes(arp.ArpPacketFromBytes(ethernet.FromBytes(buf[:n]).Payload()).Payload()).SenderMac()

	fmt.Println(senderMAC.String())
}