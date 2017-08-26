package main

import (
	"bytes"
	"fmt"
	"testing"
)

var testPayload = `{"test": "payload"}`

var testPacket = Packet{
	Header: Header{
		Guard:         0xff,
		IsReply:       false,
		SessionID:     0x57abface,
		Bytes8to13:    [6]byte{1, 2, 3, 4, 5, 6},
		PayloadType:   0x70ad,
		PayloadLength: int32(len(testPayload)),
	},
	Payload: []byte(testPayload),
}

var testBinPacket = []byte{0xff, 0x0, 0x0, 0x0, 0xce, 0xfa, 0xab, 0x57, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0xad, 0x70, 0x13, 0x0, 0x0, 0x0, 0x7b, 0x22, 0x74, 0x65, 0x73, 0x74, 0x22, 0x3a, 0x20, 0x22, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x7d}

func TestMarshalBinary(t *testing.T) {
	bin, err := testPacket.MarshalBinary()
	if err != nil {
		t.Fatalf("Packet.MarshalBinary: %+v", err)
	}

	if !bytes.Equal(bin, testBinPacket) {
		t.Fatalf("Packet.MarshalBinary didn't marshal correctly")
	}
}

func TestUnmarshalBinary(t *testing.T) {
	packet := &Packet{}
	err := packet.UnmarshalBinary(testBinPacket)
	if err != nil {
		t.Fatalf("Packet.UnmarshalBinary: %+v", err)
	}

	if !packet.Equals(&testPacket) {
		t.Fatalf("Packet.UnmarshalBinary didn't unmarshal correctly")
	}
}

func TestString(t *testing.T) {
	fmt.Println(testPacket.String())
}
