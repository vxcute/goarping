package arp

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
	ArpRequest   = 0x0001
	ArpReply     = 0x0002
	ArpIpv4      = 0x0800 
	ArpEthernet  = 0x0001
	HardwareSize = 0x6
	ProtocolSize = 0x4
)

type ArpPacket struct {
	hardwareType uint16
	protocolType uint16
	hardwareSize uint8 
	protocolSize uint8 
	opcode 		 uint16
	payload 	 []byte
}

type ArpIPv4 struct {
	senderMAC  net.HardwareAddr
	senderIP   uint32
	targetMAC  net.HardwareAddr
	targetIP   uint32
}

func NewPacket(hwtype uint16,  protype uint16, hwsize uint8, protsize uint8, opcode uint16, payload []byte) *ArpPacket {
	return &ArpPacket{
		hardwareType: hwtype,
		protocolType: protype,
		hardwareSize: hwsize,
		protocolSize: protsize,
		opcode: opcode,
		payload: payload,
	}
}

func NewARPIPv4(senderMac net.HardwareAddr, destMac net.HardwareAddr, senderIP net.IP, targetIP net.IP) *ArpIPv4 {
	return &ArpIPv4{
		senderMAC: senderMac,
		targetMAC: destMac, 
		senderIP: binary.BigEndian.Uint32(senderIP.To4()),
		targetIP: binary.BigEndian.Uint32(targetIP.To4()),
	}
}

func ArpPacketFromBytes(b []byte) *ArpPacket {
	return &ArpPacket{
		hardwareType: binary.BigEndian.Uint16(b[0:2]),
		protocolType: binary.BigEndian.Uint16(b[2:4]),
		hardwareSize: b[4],
		protocolSize: b[5],
		opcode: binary.BigEndian.Uint16(b[6:8]),
		payload: b[8:],
	}
}

func ARPIPv4FromBytes(b []byte) *ArpIPv4 {
	return &ArpIPv4{
		senderMAC: b[0:6],
		senderIP: binary.BigEndian.Uint32(b[6:10]),
		targetMAC: b[10:16],
		targetIP: binary.BigEndian.Uint32(b[16:20]),
	}
}

func (ad *ArpIPv4) SenderIP() uint32 {
	return ad.senderIP
}

func (ad *ArpIPv4) TargetIP() uint32 {
	return ad.targetIP
}

func (ad *ArpIPv4) SenderMac() net.HardwareAddr {
	return ad.senderMAC
}

func (ad *ArpIPv4) TargetMac() net.HardwareAddr {
	return ad.targetMAC
}

func (ad *ArpIPv4) Bytes() []byte {
	buf := new(bytes.Buffer)
	buf.Write(ad.senderMAC)
	binary.Write(buf, binary.BigEndian, ad.senderIP)
	buf.Write(ad.targetMAC)
	binary.Write(buf, binary.BigEndian, ad.targetIP)
	return buf.Bytes()
}

func (ap *ArpPacket) HardwareType() uint16 {
	return ap.hardwareType
}

func (ap *ArpPacket) ProtocolType() uint16 {
	return ap.protocolType
}

func (ap *ArpPacket) HardwareSize() uint8 {
	return ap.hardwareSize
}

func (ap *ArpPacket) ProtocolSize() uint8 {
	return ap.protocolSize
}

func (ap *ArpPacket) Opcode() uint16 {
	return ap.opcode
}

func (ap *ArpPacket) Payload() []byte {
	return ap.payload
}

func (ap *ArpPacket) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, ap.hardwareType)
	binary.Write(buf, binary.BigEndian, ap.protocolType)
	binary.Write(buf, binary.BigEndian, ap.hardwareSize)
	binary.Write(buf, binary.BigEndian, ap.protocolSize)
	binary.Write(buf, binary.BigEndian, ap.opcode)
	buf.Write(ap.payload)
	return buf.Bytes()
}