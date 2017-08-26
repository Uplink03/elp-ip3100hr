package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/lunixbochs/struc"
)

type Packet struct {
	Guard         byte
	IsReply       bool
	Byte3         byte
	Byte4         byte
	SessionID     int32 `struc:"little"`
	Bytes8to13    [6]byte
	PayloadType   int16 `struc:"little"`
	PayloadLength int32 `struc:"little,sizeof=Payload"`
	Payload       []byte
}

const HeaderLength = 20

func (packet *Packet) MarshalBinary() (data []byte, err error) {
	buf := &bytes.Buffer{}
	err = struc.Pack(buf, packet)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (packet *Packet) UnmarshalBinary(data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	err = struc.Unpack(buf, packet)

	return err
}

func (p1 *Packet) Equals(p2 *Packet) bool {
	return (p1.Guard == p2.Guard &&
		p1.IsReply == p2.IsReply &&
		p1.Byte3 == p2.Byte3 &&
		p1.SessionID == p2.SessionID &&
		p1.Bytes8to13 == p2.Bytes8to13 &&
		p1.PayloadType == p2.PayloadType &&
		p1.PayloadLength == p2.PayloadLength &&
		bytes.Equal(p1.Payload, p2.Payload))
}

func (packet *Packet) HasBinaryPayload() bool {
	return false
}

func (packet *Packet) String() string {
	templateTxt := `Guard: {{.Guard | printf "%x"}}
IsReply: {{.IsReply}}
Byte3: {{.Byte3}} / {{.Byte3 | printf "%02x"}}
Byte4: {{.Byte4}} / {{.Byte4 | printf "%02x"}}
SessionID: {{.SessionID}} / 0x{{.SessionID | printf "%8X"}}
Bytes8to13: {{.Bytes8to13}} / {{.Bytes8to13 | printf "%02x"}}
PayloadType: {{.PayloadType}} / {{.PayloadType | printf "%04x"}}
PayloadLength: {{.PayloadLength}}
Payload:
{{if .HasBinaryPayload}}{{.Payload}}{{else}}{{.Payload | printf "%s"}}{{end}}
`
	tmpl, err := template.New("Packet").Parse(templateTxt)
	fmt.Println(err)
	buf := &bytes.Buffer{}
	tmpl.Execute(buf, packet)

	return buf.String()
}
