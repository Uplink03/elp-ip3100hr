package main

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/lunixbochs/struc"
)

type Header struct {
	Guard         byte
	IsReply       bool
	Byte3         byte
	Byte4         byte
	SessionID     int32 `struc:"little"`
	Bytes8to13    [6]byte
	PayloadType   int16 `struc:"little"`
	PayloadLength int32 `struc:"little"`
}

type Packet struct {
	Header
	Payload []byte
}

const HeaderLength = 20

func (header *Header) MarshalBinary() (data []byte, err error) {
	buf := &bytes.Buffer{}
	err = struc.Pack(buf, header)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (packet *Packet) MarshalBinary() (data []byte, err error) {
	if packet.Header.PayloadLength != int32(len(packet.Payload)) {
		return nil, errors.New("Payload length mismatch")
	}

	headerData, err := packet.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	buf.Write(headerData)
	buf.Write(packet.Payload)

	return buf.Bytes(), nil
}

func (header *Header) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	err = struc.Unpack(buf, header)

	return err
}

func (packet *Packet) UnmarshalBinary(data []byte) (err error) {
	err = packet.Header.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	if packet.PayloadLength+HeaderLength != int32(len(data)) {
		return errors.New("Payload length mismatch")
	}

	packet.Payload = make([]byte, packet.PayloadLength)
	copy(packet.Payload, data[HeaderLength:])

	return nil
}

func (h1 *Header) Equals(h2 *Header) bool {
	return (h1.Guard == h2.Guard &&
		h1.IsReply == h2.IsReply &&
		h1.Byte3 == h2.Byte3 &&
		h1.SessionID == h2.SessionID &&
		h1.Bytes8to13 == h2.Bytes8to13 &&
		h1.PayloadType == h2.PayloadType &&
		h1.PayloadLength == h2.PayloadLength)
}

func (p1 *Packet) Equals(p2 *Packet) bool {
	return p1.Header.Equals(&p2.Header) && bytes.Equal(p1.Payload, p2.Payload)
}

func (header *Header) HasBinaryPayload() bool {
	return false // TODO find out which packet types have binary payload
}

var headerTemplate, payloadTemplate *template.Template

func init() {
	headerTemplateTxt := `Guard: {{.Guard | printf "%x"}}
IsReply: {{.IsReply}}
Byte3: {{.Byte3}} / {{.Byte3 | printf "%02x"}}
Byte4: {{.Byte4}} / {{.Byte4 | printf "%02x"}}
SessionID: {{.SessionID}} / 0x{{.SessionID | printf "%08X"}}
Bytes8to13: {{.Bytes8to13}} / {{.Bytes8to13 | printf "%02x"}}
PayloadType: {{.PayloadType}} / {{.PayloadType | printf "%04x"}}
PayloadLength: {{.PayloadLength}}
`
	payloadTemplateTxt := `Payload:
{{if .HasBinaryPayload}}{{.Payload}}{{else}}{{.Payload | printf "%s"}}{{end}}
`
	headerTemplate, _ = template.New("PacketHeader").Parse(headerTemplateTxt)
	payloadTemplate, _ = template.New("Packet").Parse(payloadTemplateTxt)
}

func (header *Header) String() string {
	buf := &bytes.Buffer{}
	headerTemplate.Execute(buf, header)

	return buf.String()
}

func (packet *Packet) String() string {
	headerTxt := packet.Header.String()

	buf := &bytes.Buffer{}
	payloadTemplate.Execute(buf, packet)
	payloadTxt := buf.String()

	return headerTxt + payloadTxt
}
