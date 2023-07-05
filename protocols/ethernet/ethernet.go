package ethernet

import (
	"bytes"
	"encoding/binary"
	"net"
)

const (
	EthernetArp uint16 = 0x0806
)

type Frame struct {
	destination net.HardwareAddr
	source		net.HardwareAddr
	etherType 	uint16
	payload 	[]byte
}

func (f *Frame) Bytes() []byte {
	return bytes.Join([][]byte{f.destination, f.source, {byte(f.etherType >> 8), byte(f.etherType)}, f.payload}, []byte{})
}

func NewFrame(dest net.HardwareAddr, src net.HardwareAddr, etherType uint16, payload []byte) *Frame {
	return &Frame{
		destination: dest,
		source: src,
		etherType: etherType,
		payload: payload,
	}
}

func FromBytes(b []byte) *Frame {
	return &Frame{
		destination: b[0:6],
		source: b[6:12],
		etherType: binary.BigEndian.Uint16(b[12:14]),
		payload: b[14:],
	}
}

func (f *Frame) Destination() net.HardwareAddr {
	return f.destination
}

func (f *Frame) Source() net.HardwareAddr {
	return f.source
}

func (f *Frame) Ethertype() uint16 {
	return f.etherType
}

func (f *Frame) Payload() []byte {
	return f.payload
}

/*
func (f *Frame) SetDestination(dest net.HardwareAddr) {
	if len(f.destination) == 0 {
		f.destination = make(net.HardwareAddr, len(dest))
	}
	
	copy(f.destination, dest)
}

func (f *Frame) SetSource(src net.HardwareAddr) {

	if len(f.source) == 0 {
		f.source = make(net.HardwareAddr, len(src))
	}

	copy(f.source, src)
}

func (f *Frame) SetEtherType(etherType uint16) {
	f.etherType = etherType
}

func (f *Frame) SetPayload(payload []byte) {
	copy(f.payload, payload)
}
*/